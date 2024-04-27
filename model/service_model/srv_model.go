package model_srv

import "gorm.io/gorm"

type Good struct {
	Description string `json:"description"`
	MonryCent   uint   `json:"monry_cent"`
}

// wx
type Payer struct {
	Openid string `json:"openid"`
}

// wx
type Amount struct {
	Total         int
	PayerTotal    int
	DiscountTotal int
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
	UsedFlage      int
}

type HttpChargeBlance struct {
	gorm.Model
	UserID        string
	Openid        string
	Phone         string
	Blance        int
	StoreID       string
	OutTradeNo    string
	TransactionId string
}

type HttpReduceBlance struct {
	gorm.Model
	UserID  string
	Openid  string
	Phone   string
	Blance  int
	StoreID string
	Remark  string
}

type IDS struct {
	IDSUserID  string `json:"user_id"`
	IDSOpenid  string `json:"open_id"`
	IDSPhone   string `json:"phone"`
	IDSStoreID string `json:"store_id"`
}

type IDSO struct {
	IDSUserID  string `json:"user_id"`
	IDSOpenid  string `json:"open_id"`
	IDSPhone   string `json:"phone"`
	IDSStoreID string `json:"store_id"`
	TrNumber   string `json:"tr_number"`
}

type IDSR struct {
	IDSUserID  string
	IDSOpenid  string
	IDSPhone   string
	IDSStoreID string
	Balance    int
}

type Response struct {
	State int    `json:"state"`
	Msg   string `json:"msg"`
	Data  bool   `json:"data"`
}
