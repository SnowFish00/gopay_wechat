package wxpay_web

import (
	gopay_api "pay/gopay_api"
	model_cfg "pay/model/config_model"
	responses "pay/response"

	"github.com/gin-gonic/gin"
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

}
