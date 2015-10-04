package main

import (
    "fmt"
    "net/http"
    "os"
    "bytes"
    "io/ioutil"
    "log"
    "errors"
    "encoding/json"
)

func main() {
    method := os.Args[1]
    params := os.Args[2:]

    reqJSON, err := makeRequest(method, params)
    if err != nil {
        log.Fatal(err)
    }

    url := "http://localhost:8080/rpc"

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqJSON))
    if err != nil {
        log.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
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

    if contents["error"] != nil {
        fmt.Println(contents["error"])
    } else {
        printResponse(method, contents)
    }
}

func printResponse(method string, contents map[string]interface{}) {
    result := contents["result"].(map[string]interface{})
    switch {
        case method == "Buy":
            printBuyResponse(result)
        case method == "Review":
            printReviewResponse(result)
        default:
            fmt.Print("Cannot happen");
    }
}

func printBuyResponse(result map[string]interface{}) {
    fmt.Print("tradeId: "); fmt.Println(result["TradeID"])
    fmt.Print("stocks: ")
    for _, stock := range result["Stocks"].([]interface{}) {
        fmt.Print(stock); fmt.Print("   ")
    }
    fmt.Println()
    fmt.Print("unvestedAmount: $"); fmt.Println(result["UnvestedAmount"])
}

func printReviewResponse(result map[string]interface{}) {
    fmt.Print("stocks: ")
    for _, stock := range result["Stocks"].([]interface{}) {
        fmt.Print(stock); fmt.Print("   ")
    }

    fmt.Println()
    fmt.Print("currentMarketValue: $"); fmt.Println(result["CurrentMarketValue"])
    fmt.Print("unvestedAmount: $"); fmt.Println(result["UnvestedAmount"])
}

func makeRequest(method string, params []string) ([]byte, error) {
    switch {
        case method == "Buy":
            if len(params) != 2 {
                return nil, errors.New("Argument is not valid")
            } else {
                return makeBuyRequest(params), nil
            }
        case method == "Review":
            if len(params) != 1 {
                return nil, errors.New("Argument is not valid")
            } else {
                return makeReviewRequest(params), nil
            }
        default:
            return nil, errors.New("ERROR: No method found with " + method)
    }
}

func makeBuyRequest(params []string) []byte {
    stockSymbolAndPercentage := params[0]
    budget := params[1]

    jsonStr := `{
                    "method":"StockDealer.Buy",
                    "params":[{
                        "StockSymbolAndPercentage": "` + stockSymbolAndPercentage + `",
                        "Budget":` + budget +
                    `}],
                    "id":0
                }`

    return []byte(jsonStr)
}

func makeReviewRequest(params []string) []byte {
    tradeID := params[0]

    jsonStr := `{
                    "method":"StockDealer.Review",
                    "params":[{
                        "TradeID":` + tradeID +
                    `}],
                    "id":0
                }`

    return []byte(jsonStr)
}
