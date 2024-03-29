package application

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/log"
	"baby-fried-rice/internal/pkg/module/file/db"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

type MinIOManager struct {
	ctx    context.Context
	lc     log.Logging
	client *minio.Client
}

func NewMinIOManager(lc log.Logging, metaId int) (FileManager, error) {
	ctx := context.Background()
	var ossMeta tables.OssMeta
	err := db.GetDB().GetObject(map[string]interface{}{"id": metaId}, &ossMeta)
	if err != nil {
		return nil, err
	}
	var manager = &MinIOManager{
		ctx: ctx,
		lc:  lc,
	}
	manager.client, err = minio.New(ossMeta.Domain, &minio.Options{
		Creds:  credentials.NewStaticV4(ossMeta.AccessKey, ossMeta.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		err = errors.Wrap(err, "failed to new minio client")
		return nil, err
	}
	err = manager.client.MakeBucket(ctx, ossMeta.Bucket, minio.MakeBucketOptions{Region: "cn-south-1", ObjectLocking: false})
	if err != nil {
		exists, _ := manager.client.BucketExists(ctx, ossMeta.Bucket)
		if !exists {
			err = errors.Wrap(err, "failed to make bucket "+ossMeta.Bucket)
			return nil, err
		}
	}
	return manager, nil
}

func (m *MinIOManager) UploadFile(key, localFilePath string) (ossMeta tables.OssMeta, downUrl string, err error) {
	var metaId = 1
	err = db.GetDB().GetObject(map[string]interface{}{"id": metaId}, &ossMeta)
	if err != nil {
		m.lc.Error(err.Error())
		return
	}
	_, err = m.client.FPutObject(m.ctx, ossMeta.Bucket, key, localFilePath, minio.PutObjectOptions{})
	if err != nil {
		m.lc.Error(err.Error())
		return
	}
	downUrl = fmt.Sprintf("http://%v/%v/%v", ossMeta.Domain, ossMeta.Bucket, key)
	return
}

func (m *MinIOManager) DeleteFile(bucket, key string) error {
	return m.client.RemoveObject(m.ctx, bucket, key, minio.RemoveObjectOptions{})
}
