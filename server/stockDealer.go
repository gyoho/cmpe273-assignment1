package main

func buyStock(stocks []Stock) Portfolio {
	for i, stock := range stocks {
		stocks[i].numStocks = int(stock.investment / stock.price)
		stocks[i].unvestment = stock.investment - (stock.price * float32(stocks[i].numStocks))
	}

	var unvestedAmount float32
	for _, stock := range stocks {
		unvestedAmount += stock.unvestment
	}

	return Portfolio{stocks, unvestedAmount}
}

func checkin(portfolio Portfolio) ([]MarketInfo, float32) {
    var symblList string
    for _, stock := range portfolio.stocks {
        symblList += ("," + stock.symbol)
    }
    symblList = symblList[1:]

    var currMarketValue float32
    currMarketInfoArr := getMarketInfoArr(symblList)
    for i, currmarketInfo := range currMarketInfoArr {
        for _, stock := range portfolio.stocks {
            if currmarketInfo.symbol == stock.symbol {

                if currmarketInfo.price > stock.price {
                    currMarketInfoArr[i].sign = "+"
                } else if currmarketInfo.price < stock.price {
                    currMarketInfoArr[i].sign = "-"
                }

                currMarketValue += currmarketInfo.price * float32(stock.numStocks)
            }
        }
    }

    return currMarketInfoArr, currMarketValue
}
