package es

import (
	"baby-fried-rice/internal/pkg/module/sync/config"
	"github.com/olivere/elastic/v7"
)

var (
	client *elastic.Client
)

func GetESClient() *elastic.Client {
	return client
}

func InitElasticSearch() (err error) {
	conf := config.GetConfig()
	client, err = elastic.NewClient(
		// 关闭sniff
		elastic.SetSniff(false),
		// 设置ES服务地址，支持多个地址
		elastic.SetURL(conf.ElasticSearch.Urls...),
		// 设置基于http base auth验证的账号和密码
		elastic.SetBasicAuth(conf.ElasticSearch.Username, conf.ElasticSearch.Password),
	)
	return
}
