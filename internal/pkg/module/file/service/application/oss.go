package application

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/log"
	"baby-fried-rice/internal/pkg/module/file/config"
	"baby-fried-rice/internal/pkg/module/file/db"
	"context"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"math/rand"
	"time"
)

type OssManager struct {
	ctx            context.Context
	lc             log.Logging
	bucketManagers map[int]*storage.BucketManager
	buckets        map[string]int
	serial         int
	r              *rand.Rand
}

func newOssBucketManager(metaId int) (*storage.BucketManager, string, error) {
	var ossMeta tables.OssMeta
	err := db.GetDB().GetObject(map[string]interface{}{"id": metaId}, &ossMeta)
	if err != nil {
		return nil, "", err
	}
	var putPolicy storage.PutPolicy
	putPolicy.Scope = ossMeta.Bucket
	mac := qbox.NewMac(ossMeta.AccessKey, ossMeta.SecretKey)
	cfg := storage.Config{}
	cfg.UseHTTPS = false
	cfg.Zone = &storage.ZoneHuadong
	cfg.UseCdnDomains = false
	bucketManager := storage.NewBucketManager(mac, &cfg)
	return bucketManager, ossMeta.Bucket, nil
}

func NewOssManager(lc log.Logging) (*OssManager, error) {
	var bucketManagers = make(map[int]*storage.BucketManager)
	var buckets = make(map[string]int)
	for i := 1; i <= config.OssMetaNum; i++ {
		bm, bucket, err := newOssBucketManager(i)
		if err != nil {
			return nil, err
		}
		buckets[bucket] = i
		bucketManagers[i] = bm
	}
	var manager = &OssManager{
		ctx:            context.Background(),
		lc:             lc,
		buckets:        buckets,
		bucketManagers: bucketManagers,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return manager, nil
}

var (
	ossManager *OssManager
)

func InitOssManager() (err error) {
	ossManager, err = NewOssManager(log.Logger)
	return
}

func GetOssManager() *OssManager {
	return ossManager
}

func (m *OssManager) UploadFile(key, localFilePath string) (ossMeta tables.OssMeta, err error) {
	var metaId = m.r.Intn(config.OssMetaNum) + 1
	err = db.GetDB().GetObject(map[string]interface{}{"id": metaId}, &ossMeta)
	if err != nil {
		m.lc.Error(err.Error())
		return
	}
	var putPolicy storage.PutPolicy
	putPolicy.Scope = ossMeta.Bucket
	mac := qbox.NewMac(ossMeta.AccessKey, ossMeta.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.UseHTTPS = false
	cfg.Zone = &storage.ZoneHuadong
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	if err = formUploader.PutFile(m.ctx, &ret, upToken, key, localFilePath, &putExtra); err != nil {
		m.lc.Error(err.Error())
		return
	}
	return
}

func (m *OssManager) getMetaIdByBucket(bucket string) int {
	return m.buckets[bucket]
}

func (m *OssManager) getManagerByBucket(bucket string) *storage.BucketManager {
	return m.bucketManagers[m.getMetaIdByBucket(bucket)]
}

func (m *OssManager) GetFile(bucket, key string) (storage.FileInfo, error) {
	bm := m.getManagerByBucket(bucket)
	return bm.Stat(bucket, key)
}

func (m *OssManager) DeleteFile(bucket, key string) error {
	bm := m.getManagerByBucket(bucket)
	return bm.Delete(bucket, key)
}
