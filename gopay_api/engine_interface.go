package gopayapi

import "github.com/go-pay/gopay/wechat/v3"

type PayAPI interface {
	//openID
	// GetOpenID() (string, error)
	//小程序支付下单
	AppletPay() (string, error)
	//小程序sigin签名
	PaySignOfApplet(*wechat.AppletParams, error)
	//H5支付下单
	H5Pay() (string, error)
	//异步内容回调
	WxPayNotify()
	//同步内容回调
	WxPayNotifySyn()
	//交易查询
	WxTestV3Quer() *wechat.QueryOrderRsp
}
