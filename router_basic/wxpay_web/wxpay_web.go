package wxpay_web

import (
	"encoding/json"
	"fmt"
	"log"
	gopay_api "pay/gopay_api"
	model_srv "pay/model/service_model"
	"pay/mysql"
	responses "pay/response"
	backgroundsyn "pay/router_basic/background_syn"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay/wechat/v3"
)

func StartOrder(c *gin.Context) {
	var good model_srv.Good
	var payer model_srv.Payer
	var pay_api gopay_api.PayAPI
	var instance gopay_api.WxPayIstance
	pay_api = instance

	if err := c.ShouldBind(&good); err != nil {
		responses.FailWithMessage(responses.ParamErr, "充值参数错误", c)
		return
	}

	if err := c.ShouldBind(&payer); err != nil {
		responses.FailWithMessage(responses.ParamErr, "openid参数错误", c)
		return
	}

	PrepayID, err := pay_api.AppletPay(c, good, payer)
	if err != nil {
		responses.FailWithMessage(responses.ServerErr, "订单下单错误", c)
		return
	}

	Parms, err := pay_api.PaySignOfApplet(PrepayID)
	if err != nil {
		responses.FailWithMessage(responses.ParamErr, "订单签发失败", c)
		return
	}

	JsonDataBytes, err := json.Marshal(Parms)
	if err != nil {
		fmt.Println("JSON marshaling failed:", err)
		responses.FailWithMessage(responses.ParamErr, "参数转化json失败", c)
		return
	}

	json := string(JsonDataBytes)

	//返回给微信支付前端方法 wx.requestPayment
	responses.OkWithDetailed(json, "订单签发完成", c)

}

// 验签
func PayNotify(c *gin.Context) {
	var pay_api gopay_api.PayAPI
	var instance gopay_api.WxPayIstance
	pay_api = instance

	resultMap, result, err := pay_api.WxPayNotify(c)

	if err == nil {
		//储存支付记录
		go PushWxPaySave(result)
		responses.OkWithDetailed(resultMap, "交易完成", c)

	} else {
		responses.FailWithMessage(responses.ParamErr, "交易检验失败", c)
	}

}

func AddNotrify(c *gin.Context) {
	PushChargeMessToPayQueue(c)
}

func ReduceNotify(c *gin.Context) {
	idsr := gopay_api.PayReduce(c)
	if idsr.Balance <= 0 {
		responses.FailWithMessage(responses.ParamErr, "停止攻击行为,ok?", c)
		return
	}
	PushReduceMessToPayQueue(idsr, c)
}

func SearchOrder(c *gin.Context) {
	var pay_api gopay_api.PayAPI
	var instance gopay_api.WxPayIstance
	pay_api = instance

	no := c.Query("TransactionId")
	params, err := pay_api.WxV3Query(no)
	if err != nil {
		responses.FailWithMessage(responses.ParamErr, "交易查询失败", c)
		return
	}

	responses.OkWithDetailed(params, "查询交易结果如下:", c)

}

func PushWxPaySave(result *wechat.V3DecryptResult) {
	//支付记录保存
	mysql.SaveWxPayDetils(result)
}

func PushChargeMessToPayQueue(c *gin.Context) {
	var idso model_srv.IDSO
	if err := c.ShouldBind(&idso); err != nil {
		responses.FailWithMessage(responses.ParamErr, "参数错误", c)
		return
	}

	order := mysql.SearchOrderTotalByOpenId(idso.TrNumber)
	if order.Amount.Total == 0 {
		responses.FailWithMessage(responses.ParamErr, "订单不存在或已被使用", c)
		return
	}

	//管理后台同步
	body := backgroundsyn.ChargeAddSyn(idso, order.Amount.Total)

	var response model_srv.Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		responses.FailWithMessage(responses.ServerErr, "管理系统数据同步失败", c)
		return
	}

	if response.State == 200 {
		//支付数据库同步
		err := mysql.BackGroundSynAdd(idso, order)
		if err != nil {
			log.Println(err.Error())
			responses.FailWithMessage(responses.SQLErr, "数据库同步失败", c)
			return
		}

		err = mysql.OrderUseOver(idso.TrNumber)
		if err != nil {
			log.Println(err.Error())
			responses.FailWithMessage(responses.SQLErr, "订单状态同步失败", c)
			return
		}

		responses.OkWithMessage("数据库同步完成", c)
	}

}

func PushReduceMessToPayQueue(idsr model_srv.IDSR, c *gin.Context) {
	//二次操作
	body := backgroundsyn.ChargeReduceSyn(idsr)

	if len(body) == 0 {
		responses.FailWithMessage(responses.ServerErr, "后台管理同步错误", c)
		return
	}

	var response model_srv.Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		responses.FailWithMessage(responses.ServerErr, "参数转化错误", c)
		return
	}

	if response.State == 200 {
		err := mysql.BackGroundSynReduce(idsr)
		if err != nil {
			responses.FailWithMessage(responses.SQLErr, "数据库同步失败", c)
			log.Println(err.Error())
			return
		}

		responses.OkWithMessage("数据库同步完成", c)
	}

}
