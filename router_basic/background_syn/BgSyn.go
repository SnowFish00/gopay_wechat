package backgroundsyn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"pay/global"
	model_srv "pay/model/service_model"
)

func ChargeAddSyn(IDSO model_srv.IDSO, balance int) []byte {
	requestData := map[string]interface{}{
		"userId":  IDSO.IDSUserID,
		"storeId": IDSO.IDSStoreID,
		"balance": balance,
	}

	jsonStr, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8848/ajax/admin/balance/recharge", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("Cookie", global.ReturnCfg().HttpServer.AdminToken)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	return body

}

func ChargeReduceSyn(idsrs model_srv.IDSRS) []byte {
	requestData := map[string]interface{}{
		"userId":  idsrs.IDSUserID,
		"storeId": idsrs.IDSStoreID,
		"balance": idsrs.Balance,
		"remark":  "系统操作扣费",
	}

	jsonStr, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8848/ajax/admin/balance/reduceBalance", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("Cookie", global.ReturnCfg().HttpServer.AdminToken)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	return body

}
