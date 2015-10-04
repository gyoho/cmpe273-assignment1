package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"strconv"
)

type MarketInfo struct {
	symbol string
	price float32
	sign string
}


func getMarketInfoArr(symblList string) []MarketInfo {
	resp, err := http.Get("http://finance.yahoo.com/webservice/v1/symbols/" + symblList + "/quote?format=json")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var contents map[string]interface{}
	err = json.Unmarshal(body, &contents)
    if err != nil {
        log.Fatal(err)
    }

	var marketInfoArr []MarketInfo
	list := contents["list"].(map[string]interface{})["resources"].([]interface{})
	for _, stock := range list {
		fields := stock.(map[string]interface{})["resource"].(map[string]interface{})["fields"]
		symbol := fields.(map[string]interface{})["symbol"].(string)
		price, err := strconv.ParseFloat(fields.(map[string]interface{})["price"].(string), 32)
		if err != nil {
			log.Fatal(err)
		}

		marketInfo := MarketInfo{symbol, float32(price), ""}
		marketInfoArr = append(marketInfoArr, marketInfo)
	}

	return marketInfoArr
}
