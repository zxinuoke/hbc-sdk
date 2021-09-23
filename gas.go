package hbc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var oracleUrl = "https://explorer.hbtcchain.io/api/v1/default_fee"

func GetHbcGas() (*HbcGas, error) {
	var result HbcGas

	err := httpGet(oracleUrl, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func httpGet(url string, obj interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, obj)
}
