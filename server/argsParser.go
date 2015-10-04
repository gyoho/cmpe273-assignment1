package main

import (
    "strings"
    "strconv"
    "log"
    "errors"
)

type Investment struct {
    symbol string
    amount float32
}


// request - "SYMBl:Pcrt"
func parseArgs(stockSymbolAndPercentage string, budget float32) ([]Investment, string, error) {
    // parse input - comma separated values
    symblPrctArray := strings.Split(stockSymbolAndPercentage, ",")

    err := checkError(symblPrctArray)
    if err != nil {
        return nil, "", err
    }

    investmentArr, symblList := constructVals(symblPrctArray, budget)

    return investmentArr, symblList, nil
}


func checkError(symblPrctArray []string) error {
    // error handling: percent should add up to 100%
    var totalPcrt int
    for _, symblPrct := range symblPrctArray {
        // get percentages  i.e) "GOOG:50%"  --> {"GOOG", "50%"}
        prctStr := strings.Split(symblPrct, ":")[1]
        prctStr = prctStr[:len(prctStr)-1]
        prct, err := strconv.Atoi(prctStr)
        if err != nil {
            log.Fatal(err)
        }
        totalPcrt += prct
    }
    if totalPcrt != 100 {
        return errors.New("ERROR: Percentages doesn't sum up to 100%")
    }
    return nil
}

func constructVals(symblPrctArray []string, budget float32) ([]Investment, string) {
    // make request string and investment objects
    var symblList string
    var investmentArr []Investment

    for _, symblPrct := range symblPrctArray {
        symbl := strings.Split(symblPrct, ":")[0]
        prctStr := strings.Split(symblPrct, ":")[1]
        prctStr = prctStr[:len(prctStr)-1]
        prct, err := strconv.Atoi(prctStr)
        if err != nil {
            log.Fatal(err)
        }
        amount := budget * float32((float32(prct)/100))

        investment := Investment{symbl, amount}
        investmentArr = append(investmentArr, investment)

        symblList += ("," + symbl)
    }

    return investmentArr, symblList[1:]    // remove the 1st comma
}
