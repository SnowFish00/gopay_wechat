package test

import (
	"context"
	"fmt"
	"pay/global"
	"pay/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
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

func TestPaysigin(c *gin.Context) {

	//获取回调参数
	notifyReq, err := wechat.V3ParseNotify(c.Request)
	if err != nil {
		fmt.Printf("回调内容异常:%v", err.Error())
		// c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "回调内容异常"})
		return
	}

	cfg := global.ReturnCfg()

	// client
	clientV3 := global.ReturnClient()

	// 验证异步通知的签名
	err = notifyReq.VerifySignByPK(clientV3.WxPublicKey())
	if err != nil {
		fmt.Printf("内容验证失败:%v", err.Error())
		// c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "内容验证失败"})
		return
	}

	// 普通支付通知解密
	result, rErr := notifyReq.DecryptCipherText(cfg.WxClient.ApiV3Key)
	if rErr != nil {
		fmt.Printf("内容解密失败:%v", err.Error())
		// c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "内容解密失败"})
		return
	}

	//success
	if result != nil && result.TradeState == "SUCCESS" {
		var wxReq = make(map[string]interface{})
		//优惠总额
		promotionAmount := 0
		for i := range result.PromotionDetail {
			promotionAmount += result.PromotionDetail[i].Amount
		}

		//商户订单号:商户系统内部订单号
		wxReq["pay_no"] = result.OutTradeNo
		//微信支付订单号:微信支付系统生成的订单号。
		wxReq["trade_no"] = result.TransactionId
		//与支付宝同步
		wxReq["trade_status"] = "TRADE_SUCCESS"
		wxReq["notify_time"] = result.SuccessTime
		//用户应支付总额
		wxReq["total_amount"] = result.Amount.Total
		//实际支付总额与优惠总额
		wxReq["receipt_amount"] = result.Amount.PayerTotal + promotionAmount
		var mapData = make(map[string]interface{})
		mapData["data_type"] = "PayNotify"
		mapData["param"] = map[string]interface{}{"payType": "wxPay", "notifyReq": wxReq}

		fmt.Println("一切就绪下面是数据")
		fmt.Println(mapData)
	}
}
