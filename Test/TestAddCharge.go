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

func BackGroundSynAddTest(IDS model_srv.IDS, Total int, OutNo string, TraID string) error {

	db := global.ReturnDB()

	toSave := model_srv.HttpChargeBlance{
		UserID:        IDS.IDSUserID,
		Openid:        IDS.IDSOpenid,
		Phone:         IDS.IDSPhone,
		Blance:        Total,
		StoreID:       IDS.IDSStoreID,
		OutTradeNo:    OutNo,
		TransactionId: TraID,
	}

	saveResult := db.Create(&toSave)

	if saveResult.Error != nil || saveResult.RowsAffected == 0 {
		return errors.New("db error")
	} else {
		fmt.Println("User created successfully")
		return nil
	}

}

func AddTest() {
	ids := model_srv.IDS{
		IDSUserID:  "5",
		IDSStoreID: "1",
		IDSPhone:   "13919898999",
	}

	body := backgroundsyn.ChargeAddSyn(ids, 100)

	var response model_srv.Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	log.Println(response.State)

	if response.State == 200 {
		err := BackGroundSynAddTest(ids, 100, "afr56&44372", "wx09#558&11")
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

}

func ReduceTest() {
	idsr := model_srv.IDSR{
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
