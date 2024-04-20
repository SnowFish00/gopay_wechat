package gopayapi

import (
	model_cfg "pay/model/config_model"
	model_srv "pay/model/service_model"

	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay/wechat/v3"
)

type PayAPI interface {
	//openID
	// GetOpenID() (string, error)
	//小程序支付下单
	AppletPay(c *gin.Context, Good model_cfg.Good) (string, error)
	//小程序sigin签名
	PaySignOfApplet(Prepayid string) (*wechat.AppletParams, error)
	//H5支付下单
	H5Pay(c *gin.Context, Good model_cfg.Good) (*wechat.H5Rsp, error)
	//异步内容回调
	WxPayNotify(c *gin.Context) (model_srv.IDS, map[string]interface{}, *wechat.V3DecryptResult, error)
	//同步内容回调
	WxPayNotifySyn(wxRsp *wechat.PrepayRsp)
	//交易查询
	WxV3Query(no string) *wechat.QueryOrderRsp
}
