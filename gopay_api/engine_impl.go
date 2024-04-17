package gopayapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	model_cfg "pay/model/config_model"
	"pay/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
)

const (
	code2sessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

func GetOpenID(c *gin.Context, wx_cfg model_cfg.WxClient) (string, error) {
	//获取code
	code := c.PostForm("code")

	//调用auth.code2Session接口获取openid
	url := fmt.Sprintf(code2sessionURL, wx_cfg.AppID, wx_cfg.AppSecret, code)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	var wxMap map[string]string
	err = json.NewDecoder(resp.Body).Decode(&wxMap)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return wxMap["openid"], err
}

func AppletPay(c *gin.Context, Good model_cfg.Good, cfg model_cfg.Config) (PrepayID string, Error error) {
	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	uuid := utils.NewUuid32
	bm := make(gopay.BodyMap)
	openid, err := GetOpenID(c, cfg.WxClient)
	if err != nil {
		return "", err
	}

	bm.Set("appid", cfg.WxClient.AppID).
		Set("mchid", cfg.WxClient.MchId).
		Set("description", Good.Description).
		Set("out_trade_no", uuid).
		Set("time_expire", expire).
		Set("notify_url", cfg.WxClient.NotifyUrl).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", Good.Monry_cent).
				Set("currency", "CNY")
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", openid)
		})

	clientV3 := cfg.NewClientV3Engine()

	wxRsp, err := clientV3.V3TransactionJsapi(context.Background(), bm)
	if err != nil {
		return "", err
	}

	return wxRsp.Response.PrepayId, errors.New(wxRsp.Error)

}

func H5Pay(c *gin.Context, Good model_cfg.Good, cfg model_cfg.Config) (wxRsp *wechat.H5Rsp, Error error) {
	uuid := utils.NewUuid32()
	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	bm := make(gopay.BodyMap)
	// openid, err := GetOpenID(c, cfg.WxClient)
	// if err != nil {
	// 	return "", err
	// }
	bm.Set("appid", cfg.WxClient.AppID).
		Set("mchid", cfg.WxClient.MchId).
		Set("description", Good.Description).
		Set("out_trade_no", uuid).
		Set("time_expire", expire).
		Set("notify_url", cfg.WxClient.NotifyUrl).
		SetBodyMap("amount", func(b gopay.BodyMap) {
			b.Set("total", Good.Monry_cent).
				Set("currency", "CNY")
		}).
		SetBodyMap("scene_info", func(b gopay.BodyMap) {
			b.Set("payer_client_ip", c.ClientIP()).
				//终端设备号
				Set("device_id", "default").
				SetBodyMap("store_info", func(b gopay.BodyMap) {
					//商城ID
					b.Set("id", "xxx").
						Set("name", "沈阳立直麻将").
						Set("area_code", "XXXXXX").
						Set("address", "辽宁省沈阳市XXX区XXX")
				}).
				SetBodyMap("h5_info", func(b gopay.BodyMap) {
					b.Set("type", "Wap").
						Set("app_name", "XXX").
						Set("app_url", "https://xxx.com").
						Set("bundle_id", "xxxxxx")
				})
		})

	clientV3 := cfg.NewClientV3Engine()

	wxRsp, err := clientV3.V3TransactionH5(context.Background(), bm)
	if err != nil {
		return nil, err
	}

	return wxRsp, errors.New(wxRsp.Error)
}

func PaySignOfApplet(cfg model_cfg.Config, Prepayid string) (*wechat.AppletParams, error) {
	clientV3 := cfg.NewClientV3Engine()
	// 小程序
	applet, err := clientV3.PaySignOfApplet("appid", "prepayid")
	return applet, err
}

/*
该函数为微信开发者工具语法函数
wx.requestPayment({
	timeStamp: '时间戳',  // 时间戳，自1970年1月1日以来的秒数
	nonceStr: '随机字符串', // 随机字符串，不长于32位
	package: 'prepay_id=xxxxx', // 统一下单接口返回的 prepay_id 参数值，提交格式如：prepay_id=*
	signType: 'MD5', // 签名算法，暂支持 MD5
	paySign: '支付签名', // 签名，具体签名方案参见微信支付开发文档
	success: function (res) {
	  // 支付成功回调
	  //dosomething
	  //......
	  console.log("支付成功", res);
	},
	fail: function (res) {
	  // 支付失败回调
	  //dosomething
	  //......
	  console.log("支付失败", res);
	}
  });
*/

func WxPayNotify(c *gin.Context, cfg model_cfg.Config) {
	notifyReq, err := wechat.V3ParseNotify(c.Request)
	if err != nil {
		c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "回调内容异常"})
		return
	}

	// client
	clientV3 := cfg.NewClientV3Engine()

	// 验证异步通知的签名
	err = notifyReq.VerifySignByPK(clientV3.WxPublicKey())
	if err != nil {
		c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "内容验证失败"})
		return
	}

	// 普通支付通知解密
	result, rErr := notifyReq.DecryptCipherText(cfg.WxClient.ApiV3Key)
	if rErr != nil {
		c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "内容解密失败"})
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
		reqJson, _ := json.Marshal(mapData)

		//可在此处向外发送reqJson
		//PushMessToPayQueue
		utils.PushMessToPayQueue(reqJson)

		c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"})

		return
	}

}

// WxTestV3Query 交易查询
func WxTestV3Query(no string, cfg model_cfg.Config) *wechat.QueryOrderRsp {
	clientV3 := cfg.NewClientV3Engine()
	wxRsp, err := clientV3.V3TransactionQueryOrder(context.Background(), wechat.OutTradeNo, no)
	if err != nil {
		return nil
	}
	return wxRsp
}
