package wxpay_web

import (
	"encoding/json"
	"fmt"
	"log"
	gopay_api "pay/gopay_api"
	model_cfg "pay/model/config_model"
	model_srv "pay/model/service_model"
	"pay/mysql"
	responses "pay/response"
	backgroundsyn "pay/router_basic/background_syn"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
)

func StartOrder(c *gin.Context) {
	var good model_cfg.Good
	var pay_api gopay_api.PayAPI
	var instance gopay_api.WxPayIstance
	pay_api = instance

	if err := c.ShouldBind(&good); err != nil {
		responses.FailWithMessage(responses.ParamErr, err.Error(), c)
		return
	}

	PrepayID, err := pay_api.AppletPay(c, good)
	if err != nil {
		responses.FailWithMessage(responses.ParamErr, err.Error(), c)
		return
	}

	Parms, err := pay_api.PaySignOfApplet(PrepayID)
	if err != nil {
		responses.FailWithMessage(responses.ParamErr, err.Error(), c)
		return
	}

	JsonDataBytes, err := json.Marshal(Parms)
	if err != nil {
		fmt.Println("JSON marshaling failed:", err)
		return
	}

	json := string(JsonDataBytes)

	//返回给微信支付前端方法 wx.requestPayment

	responses.OkWithData(json, c)

}

// 验签
func PayNotify(c *gin.Context) {
	var pay_api gopay_api.PayAPI
	var instance gopay_api.WxPayIstance
	pay_api = instance

	ids, resultMap, result, err := pay_api.WxPayNotify(c)

	if err == nil {
		//协程二次操作
		go PushMessToPayQueue(ids, resultMap, result)
		responses.OkWithDetailed(&wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"}, "ok", c)

	} else {
		responses.FailWithMessage(responses.ParamErr, "交易检验失败", c)
	}

}

func SearchOrder(c *gin.Context) {
	var pay_api gopay_api.PayAPI
	var instance gopay_api.WxPayIstance
	pay_api = instance

	no := c.Query("outTradeNo")
	params := pay_api.WxV3Query(no)

	responses.OkWithDetailed(params, "查询交易结果如下:", c)

}

func PushMessToPayQueue(ids model_srv.IDS, resultMap map[string]interface{}, result *wechat.V3DecryptResult) {
	//二次操作
	body := backgroundsyn.ChargeAddSyn(ids, result.Amount.Total)

	var response model_srv.Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	if response.State == 200 {
		err := mysql.BackGroundSynAdd(ids, result)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	log.Fatalln("http error")

}
