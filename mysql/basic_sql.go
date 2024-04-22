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
			Currency:      result.Amount.Currency,
			PayerCurrency: result.Amount.PayerCurrency,
		},
		SceneInfo: model_srv.SceneInfo{
			DeviceID: result.SceneInfo.DeviceId,
		},
	}

	saveResult := db.Create(&toSave)

	if saveResult.Error != nil || saveResult.RowsAffected == 0 {
		fmt.Println("Failed to create user:", saveResult.Error)
	} else {
		fmt.Println("User created successfully")
	}

}

func BackGroundSynAdd(IDS model_srv.IDS, result *wechat.V3DecryptResult) error {
	db := global.ReturnDB()

	toSave := model_srv.HttpChargeBlance{
		UserID:        IDS.IDSUserID,
		Openid:        IDS.IDSOpenid,
		Phone:         IDS.IDSPhone,
		Blance:        result.Amount.PayerTotal,
		StoreID:       IDS.IDSStoreID,
		OutTradeNo:    result.OutTradeNo,
		TransactionId: result.TransactionId,
	}

	saveResult := db.Create(&toSave)

	if saveResult.Error != nil || saveResult.RowsAffected == 0 {
		return errors.New("db error")
	} else {
		fmt.Println("User created successfully")
		return nil
	}

}

func BackGroundSynReduce(IDSR model_srv.IDSR) error {
	db := global.ReturnDB()

	toSave := model_srv.HttpReduceBlance{
		UserID:  IDSR.IDSUserID,
		Openid:  IDSR.IDSOpenid,
		Phone:   IDSR.IDSPhone,
		Blance:  IDSR.Balance,
		StoreID: IDSR.IDSStoreID,
	}

	saveResult := db.Where("open_id = ?", toSave.Openid).Save(&toSave)

	if saveResult.Error != nil || saveResult.RowsAffected == 0 {
		return errors.New("db error")
	} else {
		fmt.Println("User created successfully")
		return nil
	}

}
