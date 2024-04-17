package gopayapi

type PayAPI interface {
	//小程序支付下单
	AppletPay() (PrepayID string, Error error)
	//H5支付下单
	H5Pay() (PrepayID string, Error error)
}
