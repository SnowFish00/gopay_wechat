package model_srv

import "gorm.io/gorm"

// wx
type Payer struct {
	Openid string
}

// wx
type Amount struct {
	Total         int
	PayerTotal    int
	Currency      string
	PayerCurrency string
}

// wx
type SceneInfo struct {
	DeviceID string
}

// 微信支付消息
type ChargeMessage struct {
	gorm.Model
	AppId          string
	MchId          string
	OutTradeNo     string
	TransactionId  string
	TradeType      string
	TradeState     string
	TradeStateDesc string
	BankType       string
	Attach         string
	SuccessTime    string
	Payer          Payer     `gorm:"embedded"`
	Amount         Amount    `gorm:"embedded"`
	SceneInfo      SceneInfo `gorm:"embedded"`
}

type HttpChargeBlance struct {
	gorm.Model
	UserID  string
	Openid  string
	Phone   string
	Blance  string
	StoreID string
}

type HttpReduceBlance struct {
	gorm.Model
	UserID  string
	Openid  string
	Phone   string
	Blance  string
	StoreID string
	Remark  string
}
