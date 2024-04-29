package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"pay/global"
	model_srv "pay/model/service_model"
	"pay/mysql"
	backgroundsyn "pay/router_basic/background_syn"
)

func BackGroundSynAddTest(IDSO model_srv.IDSO, Total int, OutNo string, TraID string) error {

	db := global.ReturnDB()

	toSave := model_srv.HttpChargeBlance{
		UserID:        IDSO.IDSUserID,
		Openid:        IDSO.IDSOpenid,
		Phone:         IDSO.IDSPhone,
		Blance:        Total,
		StoreID:       IDSO.IDSStoreID,
		OutTradeNo:    OutNo,
		TransactionId: TraID,
	}

	saveResult := db.Create(&toSave)

	if saveResult.Error != nil || saveResult.RowsAffected == 0 {
		return errors.New("db error")
	} else {
		fmt.Println("充值成功")
		return nil
	}

}

func AddTest() {
	idso := model_srv.IDSO{
		IDSUserID:  "5",
		IDSStoreID: "1",
		IDSPhone:   "13919898999",
	}

	body := backgroundsyn.ChargeAddSyn(idso, 100)

	var response model_srv.Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	log.Println(response.State)

	if response.State == 200 {
		err := BackGroundSynAddTest(idso, 100, "afr56&44372", "wx09#558&11")
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

}

func ReduceTest() {
	idsr := model_srv.IDSRS{
		IDSUserID:  "5",
		IDSStoreID: "1",
		IDSPhone:   "13919898999",
		Balance:    100,
	}

	body := backgroundsyn.ChargeReduceSyn(idsr)

	var response model_srv.Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	log.Println(response.State)

	if response.State == 200 {
		err := mysql.BackGroundSynReduce(idsr)
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

}
