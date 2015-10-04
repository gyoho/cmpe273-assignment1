package main

import (
	"math/rand"
	"errors"
)

type Portfolio struct {
    stocks []Stock
    unvestedAmount float32
}

type Stock struct {
	symbol string
	price float32
	numStocks int
	investment float32
	unvestment float32
}

func deal(stockSymbolAndPercentage string, budget float32) (float32, Portfolio, error) {
	investmentArr, symblList, err := parseArgs(stockSymbolAndPercentage, budget)
	if err != nil {
		return 0.0, Portfolio{}, err
	}

	var stocks []Stock
	for _, investment := range investmentArr {
		stock := Stock{investment.symbol, 0.0, 0, investment.amount, 0.0}
		stocks = append(stocks, stock)
	}

	marketInfoArr := getMarketInfoArr(symblList)
	for _, marketInfo := range marketInfoArr {
		for i, stock := range stocks {
			if marketInfo.symbol == stock.symbol {
				stocks[i].price = marketInfo.price
			}
		}
	}


	portfolio := buyStock(stocks)

	tradeId := rand.Float32()
	TradeIDs = append(TradeIDs, tradeId)
	Portfolios[tradeId] = portfolio

	return tradeId, portfolio, nil
}

func checkPortfolio(tradeId float32) (Portfolio, []MarketInfo, float32, error) {
	portfolio, ok := Portfolios[tradeId]
	if ok {
		currMarketInfoArr, currMarketValue := checkin(portfolio)
		return portfolio, currMarketInfoArr, currMarketValue, nil
	} else {
		return Portfolio{}, nil, 0.0, errors.New("ERROR: No portfolio found with this ID")
	}
}
