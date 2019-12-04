package service

import (
	"testing"
	"github.com/3115826227/baby-fried-rice/module/crawler/model"
	"fmt"
)

func TestStation(t *testing.T) {
	Station()
}

func TestTrainsMeta(t *testing.T) {
	//TrainsMeta("K7", "20191201")
}

func TestGetTrainWay2(t *testing.T) {
	//GetTrainWay2("K73")
}

func TestExecutor(t *testing.T) {
	//Executor("2019-12-01")
}

func TestTrainsMetaExtract(t *testing.T) {
	TrainMetaExecutor("2019-12-10")
}

func TestMeituanRun(t *testing.T) {
	var train = model.TrainMeta{
		TrainCode:     "",
		StartStation:  "金华",
		ArriveStation: "娄底",
	}
	fmt.Println(train)
}

func TestGetStation(t *testing.T) {
	list := GetStation()
	fmt.Println(list)
}

func TestUpdateCity(t *testing.T) {
	UpdateCity("2019-12-04")
}

func TestTrainSeatPrice(t *testing.T) {
	TrainSeatPrice("2019-12-02")
}

func TestCrawler(t *testing.T) {
	Crawler(SeatRelationInfo{Date: "2019-12-16", From: "北京", To: "衡阳"})
}

func TestTongChengYiLongTraffic(t *testing.T) {
	//TongChengYiLongTraffic(SeatRelationInfo{Date: "2019-12-16", From: "北京", To: "衡阳"})
}

func TestMeituanTraffic(t *testing.T) {
	//MeituanTraffic(SeatRelationInfo{Date: "2019-12-16", From: "北京", To: "衡阳"})
}

func TestJindongTraffic(t *testing.T) {
	//JindongTraffic(SeatRelationInfo{Date: "2019-12-16", From: "北京", To: "衡阳"})
}

func TestZhixingTrainTrigger(t *testing.T) {
	//ZhixingTrainTrigger(TrainTask{Train: model.TrainMeta{Train: "T109", StartStation: "北京", ArriveStation: "上海", Date: "2019-12-10"}, Retry: 0})
}

func TestTrainMetaTrigger(t *testing.T) {
	//MeituanTrainTrigger(TrainTask{Train: model.TrainMeta{Train: "T109", StartStation: "北京", ArriveStation: "上海", Date: "2019-12-10"}, Retry: 0})
}
