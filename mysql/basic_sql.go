package mysql

import (
	"errors"
	"fmt"
	"pay/global"
	model_srv "pay/model/service_model"

	"github.com/go-pay/gopay/wechat/v3"
)

func SaveWxPayDetils(result *wechat.V3DecryptResult) {
	db := global.ReturnDB()
	toSave := model_srv.ChargeMessage{
		AppId:          result.Appid,
		MchId:          result.Mchid,
		OutTradeNo:     result.OutTradeNo,
		TransactionId:  result.TransactionId,
		TradeType:      result.TradeType,
		TradeState:     result.TradeState,
		TradeStateDesc: result.TradeStateDesc,
		BankType:       result.BankType,
		Attach:         result.Attach,
		SuccessTime:    result.SuccessTime,
		Payer: model_srv.Payer{
			Openid: result.Payer.Openid,
		},
		Amount: model_srv.Amount{
			Total:         result.Amount.Total,
			PayerTotal:    result.Amount.PayerTotal,
			DiscountTotal: result.Amount.DiscountTotal,
			Currency:      result.Amount.Currency,
			PayerCurrency: result.Amount.PayerCurrency,
		},
		SceneInfo: model_srv.SceneInfo{
			DeviceID: "normal",
		},
	}

	saveResult := db.Create(&toSave)

	if saveResult.Error != nil || saveResult.RowsAffected == 0 {
		fmt.Printf("支付记录%s保存失败:\n", result.TransactionId)
	} else {
		fmt.Printf("支付记录%s保存成功:\n", result.TransactionId)
	}

}

func BackGroundSynAdd(IDSO model_srv.IDSO, result model_srv.ChargeMessage) error {
	db := global.ReturnDB()

	toSave := model_srv.HttpChargeBlance{
		UserID:        IDSO.IDSUserID,
		Openid:        IDSO.IDSOpenid,
		Phone:         IDSO.IDSPhone,
		Blance:        result.Amount.Total,
		StoreID:       IDSO.IDSStoreID,
		OutTradeNo:    result.OutTradeNo,
		TransactionId: result.TransactionId,
	}

	saveResult := db.Create(&toSave)

	if saveResult.Error != nil || saveResult.RowsAffected == 0 {
		return errors.New("db error")
	} else {
		fmt.Println("充值成功")
		return nil
	}

}

func BackGroundSynReduce(IDSRS model_srv.IDSRS) error {
	db := global.ReturnDB()

	toSave := model_srv.HttpReduceBlance{
		UserID:  IDSRS.IDSUserID,
		Openid:  IDSRS.IDSOpenid,
		Phone:   IDSRS.IDSPhone,
		Blance:  IDSRS.Balance,
		StoreID: IDSRS.IDSStoreID,
	}

	saveResult := db.Where("open_id = ?", toSave.Openid).Save(&toSave)

	if saveResult.Error != nil || saveResult.RowsAffected == 0 {
		return errors.New("db error")
	} else {
		fmt.Println("充值成功")
		return nil
	}

}

func SearchOrderTotalByOpenId(Trno string) model_srv.ChargeMessage {
	db := global.ReturnDB()
	var order model_srv.ChargeMessage
	saveResult := db.Where("transaction_id = ?", Trno).Find(&order)
	if saveResult.Error != nil || saveResult.RowsAffected == 0 {
		fmt.Println("查找支付失败")
		return model_srv.ChargeMessage{}
	}
	if order.UsedFlage == 1 {
		fmt.Println("订单已被记录")
		return model_srv.ChargeMessage{}
	}
	return order

}

func OrderUseOver(Trno string) error {
	db := global.ReturnDB()
	var order model_srv.ChargeMessage
	searchResult := db.Where("transaction_id =?", Trno).First(&order)
	if searchResult.Error != nil || searchResult.RowsAffected == 0 {
		fmt.Println("未找到订单")
		return searchResult.Error
	}

	order.UsedFlage = 1

	saveResult := db.Save(&order)

	if saveResult.Error != nil || saveResult.RowsAffected == 0 {
		fmt.Println("修改失败")
		return saveResult.Error
	}

	return nil
}
