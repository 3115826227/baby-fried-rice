package application

import (
	"baby-fried-rice/internal/pkg/kit/db/tables"
	"baby-fried-rice/internal/pkg/kit/rpc/pbservices/sms"
	"baby-fried-rice/internal/pkg/module/smsDao/config"
	"baby-fried-rice/internal/pkg/module/smsDao/db"
	"baby-fried-rice/internal/pkg/module/smsDao/log"
	"context"
	"encoding/json"
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"time"
)

type SmsService struct {
}

var (
	client *dysmsapi20170525.Client
)

func InitClient() (err error) {
	var c = config.GetConfig()
	cfg := &openapi.Config{
		AccessKeyId:     &c.AccessKeyId,
		AccessKeySecret: &c.AccessKeySecret,
	}
	cfg.Endpoint = tea.String(c.Endpoint)
	client, err = dysmsapi20170525.NewClient(cfg)
	return
}

type TemplateSmsParam struct {
	Code *string `json:"code"`
}

func SendMessage(phone *string, param string) (err error) {
	if param == "" {
		err = errors.New("param must no empty")
		return
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("宝宝煎米果平台"),
		PhoneNumbers:  phone,
		TemplateCode:  tea.String("SMS_164100177"),
		TemplateParam: tea.String(param),
	}
	// 复制代码运行请自行打印 API 的返回值
	resp, err := client.SendSms(sendSmsRequest)
	if err != nil {
		return err
	}
	log.Logger.Debug(resp.String())
	return err
}

func (service *SmsService) SendMessageDao(ctx context.Context, req *sms.ReqSendMessageDao) (empty *emptypb.Empty, err error) {
	// 查找模板
	var smt tables.SendMessageTemplate
	if err = db.GetDB().GetObject(map[string]interface{}{"code": req.Code}, &smt); err != nil {
		// 找不到对应模板
		if err == gorm.ErrRecordNotFound {
			err = errors.New("找不到对应的短信模板")
			log.Logger.Error(err.Error())
			return
		}
		log.Logger.Error(err.Error())
		return
	}
	// 查找发送记录，控制单个手机号发送频率不能超过1分钟/条
	var latestLog tables.SendMessageLog
	if err = db.GetDB().GetDB().Where("phone = ?", req.Phone).
		Order("send_timestamp").First(&latestLog).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Logger.Error(err.Error())
			return
		}
		// 没有发送记录
		err = nil
	} else {
		// 当前时间的前一分钟小于最近一条发送时间
		if time.Now().Add(-time.Minute).Unix() < latestLog.SendTimestamp {
			err = errors.New("您发送的频率太频繁，请稍后再试")
			log.Logger.Error(err.Error())
			return
		}
	}
	var l = tables.SendMessageLog{
		AccountId:     req.AccountId,
		Phone:         req.Phone,
		Code:          req.PhoneCode,
		SignName:      req.SignName,
		TemplateCode:  smt.TemplateCode,
		SendTimestamp: time.Now().Unix(),
	}
	var tx = db.GetDB().GetDB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()
	// 创建日志
	if err = tx.Create(&l).Error; err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 配置参数
	var mp = map[string]interface{}{
		"code": req.PhoneCode,
	}
	var data []byte
	if data, err = json.Marshal(mp); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	// 发送短信
	if err = SendMessage(tea.String(req.Phone), string(data)); err != nil {
		log.Logger.Error(err.Error())
		return
	}
	return
}
