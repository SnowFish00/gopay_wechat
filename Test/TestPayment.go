package test

import (
	"context"
	"fmt"
	"pay/global"
	"pay/utils"
	"time"

	"github.com/go-pay/gopay"
)

func TestAppltPay() (PrepayID string, Error error) {
	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)

	uuid := utils.NewUuid32()

	bm := make(gopay.BodyMap)

	cfg := global.ReturnCfg()

	bm.Set("appid", cfg.WxClient.AppID).
		Set("mchid", cfg.WxClient.MchId).
		Set("description", "一分钱测试").
		Set("out_trade_no", uuid).
		Set("time_expire", expire).
		Set("notify_url", cfg.WxClient.NotifyUrl).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", 1).
				Set("currency", "CNY")
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", cfg.Test.OpenID)
		})

	fmt.Println("构造如下")
	fmt.Println(bm)

	clientV3 := global.ReturnClient()

	wxRsp, err := clientV3.V3TransactionJsapi(context.Background(), bm)
	if err != nil {
		return "", err
	}

	return wxRsp.Response.PrepayId, err
}
