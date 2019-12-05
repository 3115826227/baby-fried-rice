package db

import (
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/crawler/config"
	"github.com/3115826227/baby-fried-rice/module/crawler/log"
	"github.com/3115826227/baby-fried-rice/module/crawler/model"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	DB *gorm.DB
)

func init() {
	var err error
	DB, err = gorm.Open("postgres", config.Config.PostgresUrl)
	fmt.Println(config.Config.PostgresUrl)
	if err != nil {
		panic(err)
	} else {
		DB.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false)
		DB.SingularTable(true)
	}
	Sync(DB)
}

func Sync(engine *gorm.DB) {
	err := engine.AutoMigrate(
		new(model.Station),
		new(model.TrainMeta),
		new(model.TrainStationRelation),
		new(model.TrainStationSeatPrice),
		new(model.TrainSeatCategory),
		new(model.TrainCategory),
	).Error
	if err != nil {
		log.Logger.Warn(err.Error())
	}
}

func AddTrainSeatCategory() {
	var categoryList = []model.TrainSeatCategory{
		{Name: "G:二等座"},
		{Name: "G:一等座"},
		{Name: "G:商务座"},
		{Name: "D:二等座"},
		{Name: "D:一等座"},
		{Name: "D:商务座"},
		{Name: "无座"},
		{Name: "硬座"},
		{Name: "硬卧"},
		{Name: "软卧"},
	}
	for _, item := range categoryList {
		if err := DB.Create(&item).Error; err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
	}
}

func AddTrainCategory() {
	var list = []model.TrainCategory{
		{Name: "G"},
		{Name: "C"},
		{Name: "D"},
		{Name: "Z"},
		{Name: "T"},
		{Name: "K"},
	}
	for _, item := range list {
		if err := DB.Create(&item).Error; err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
	}
}
