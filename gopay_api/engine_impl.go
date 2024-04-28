package gopayapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"pay/global"
	model_srv "pay/model/service_model"
	responses "pay/response"
	"pay/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/go-pay/xlog"
)

const (
	code2sessionURL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

type WxPayIstance struct {
}

func GetOpenIDBycode2Session(c *gin.Context) (string, error) {
	//获取code
	code := c.PostForm("code")
	wx_cfg := global.ReturnCfg()
	//调用auth.code2Session接口获取openid
	url := fmt.Sprintf(code2sessionURL, wx_cfg.WxClient.AppID, wx_cfg.WxClient.AppSecret, code)
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

func GetOpenIDByFront() string {
	return ""
}

func (WxPayIstance) AppletPay(c *gin.Context, Good model_srv.Good, Payer model_srv.Payer) (PrepayID string, Error error) {
	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)

	uuid := utils.NewUuid32()

	bm := make(gopay.BodyMap)

	openid := Payer.Openid

	cfg := global.ReturnCfg()

	bm.Set("appid", cfg.WxClient.AppID).
		Set("mchid", cfg.WxClient.MchId).
		Set("description", Good.Description).
		Set("out_trade_no", uuid).
		Set("time_expire", expire).
		Set("notify_url", cfg.WxClient.NotifyUrl).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", Good.MonryCent).
				Set("currency", "CNY")
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", openid)
		})

	clientV3 := global.ReturnClient()

	wxRsp, err := clientV3.V3TransactionJsapi(context.Background(), bm)
	if err != nil {
		return "", err
	}

	return wxRsp.Response.PrepayId, err

}

func (WxPayIstance) H5Pay(c *gin.Context, Good model_srv.Good) (wxRsp *wechat.H5Rsp, Error error) {
	uuid := utils.NewUuid32()
	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	bm := make(gopay.BodyMap)
	cfg := global.ReturnCfg()
	bm.Set("appid", cfg.WxClient.AppID).
		Set("mchid", cfg.WxClient.MchId).
		Set("description", Good.Description).
		Set("out_trade_no", uuid).
		Set("time_expire", expire).
		Set("notify_url", cfg.WxClient.NotifyUrl).
		SetBodyMap("amount", func(b gopay.BodyMap) {
			b.Set("total", Good.MonryCent).
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

	clientV3 := global.ReturnClient()

	wxRsp, err := clientV3.V3TransactionH5(context.Background(), bm)
	if err != nil {
		return nil, err
	}

	return wxRsp, err
}

func (WxPayIstance) PaySignOfApplet(Prepayid string) (*wechat.AppletParams, error) {
	cfg := global.ReturnCfg()
	clientV3 := global.ReturnClient()
	// 小程序
	applet, err := clientV3.PaySignOfApplet(cfg.WxClient.AppID, Prepayid)
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

// 异步回调 前端附带 openid phone userid
func (WxPayIstance) WxPayNotify(c *gin.Context) (map[string]interface{}, *wechat.V3DecryptResult, error) {

	notifyReq, err := wechat.V3ParseNotify(c.Request)
	if err != nil {
		// c.JSON(http.StatusBadRequest, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "回调内容异常"})
		log.Printf("回调内容异常: %v \n", err)
		return nil, nil, err
	}

	cfg := global.ReturnCfg()

	// client
	clientV3 := global.ReturnClient()

	// 验证异步通知的签名
	err = notifyReq.VerifySignByPK(clientV3.WxPublicKey())
	if err != nil {
		// c.JSON(http.StatusBadRequest, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "内容验证失败"})
		log.Printf("内容验证失败: %v \n", err)
		return nil, nil, err
	}

	// 普通支付通知解密
	result, rErr := notifyReq.DecryptCipherText(cfg.WxClient.ApiV3Key)
	if rErr != nil {
		// c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "内容解密失败"})
		log.Printf("内容解密失败: %v \n", err)
		return nil, nil, err
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

		return mapData, result, nil
	}

	return nil, nil, errors.New("FialedPayment")

}

// 同步验签
func (WxPayIstance) WxPayNotifySyn(wxRsp *wechat.PrepayRsp) {
	clientV3 := global.ReturnClient()
	pkMap := clientV3.WxPublicKeyMap()
	// wxPublicKey：微信平台证书公钥内容，通过 client.WxPublicKeyMap() 获取，然后根据 signInfo.HeaderSerial 获取相应的公钥
	err := wechat.V3VerifySignByPK(wxRsp.SignInfo.HeaderTimestamp, wxRsp.SignInfo.HeaderNonce, wxRsp.SignInfo.SignBody, wxRsp.SignInfo.HeaderSignature, pkMap[wxRsp.SignInfo.HeaderSerial])
	if err != nil {
		xlog.Error(err)
		return
	}
}

// WxTestV3Query 交易查询 no TransactionId or outTradeNo
func (WxPayIstance) WxV3Query(no string) (*wechat.QueryOrderRsp, error) {
	clientV3 := global.ReturnClient()
	wxRsp, err := clientV3.V3TransactionQueryOrder(context.Background(), wechat.TransactionId, no)
	if err != nil {
		return nil, err
	}
	return wxRsp, nil
}

// 扣费
func PayReduce(c *gin.Context) model_srv.IDSR {
	var idsr model_srv.IDSR
	if err := c.ShouldBind(&idsr); err != nil {
		responses.FailWithMessage(responses.ParamErr, "参数错误", c)
		return idsr
	}

	return idsr
}
