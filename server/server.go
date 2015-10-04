package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"fmt"

	"github.com/bakins/net-http-recover"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/justinas/alice"
)

// modified versions of examples found at http://golang.org/pkg/net/rpc/
type BuyArgs struct {
	StockSymbolAndPercentage string
    Budget float32
}

type ReviewArgs struct {
	TradeID float32
}

type Purchase struct {
	TradeID float32
    Stocks []string
    UnvestedAmount string
}

type Book struct {
	Stocks []string
	CurrentMarketValue float32
	UnvestedAmount string
}

type StockDealer struct {}


var TradeIDs []float32
var Portfolios map[float32]Portfolio


func (t *StockDealer) Buy(req *http.Request, args *BuyArgs, res *Purchase) error {

    tradeID, portfolio, err := deal(args.StockSymbolAndPercentage, args.Budget)
    if err != nil {
        return err
    }

    var stocksInfo []string
    for _, stock := range portfolio.stocks {
        // “GOOG:100:$500.25”, “YHOO:200:$31.40”
        stocksInfo = append(stocksInfo, stock.symbol + ":" + strconv.Itoa(stock.numStocks) + ":$" + strconv.FormatFloat(float64(stock.price), 'f', 2, 32))
    }

    res.TradeID = tradeID
    res.Stocks = stocksInfo
    res.UnvestedAmount = strconv.FormatFloat(float64(portfolio.unvestedAmount), 'f', 2, 32)

	return nil
}

func (t *StockDealer) Review(req *http.Request, args *ReviewArgs, res *Book) error {
	portfolio, currMarketInfoArr, currMarketValue, err := checkPortfolio(args.TradeID)
	if err != nil {
		return err
	}

	var currStocksInfo []string
    for _, stock := range portfolio.stocks {
		for _, currMarketInfo := range currMarketInfoArr {
			if stock.symbol == currMarketInfo.symbol {
				// “GOOG:100:+$520.25”, “YHOO:200:-$30.40”
				currStocksInfo = append(currStocksInfo, stock.symbol + ":" + strconv.Itoa(stock.numStocks) + ":" + currMarketInfo.sign + "$" + strconv.FormatFloat(float64(currMarketInfo.price), 'f', 2, 32))
			}
		}
    }

	res.Stocks = currStocksInfo
	res.CurrentMarketValue = currMarketValue
	res.UnvestedAmount = strconv.FormatFloat(float64(portfolio.unvestedAmount), 'f', 2, 32)

	return nil
}

func main() {
	Portfolios = make(map[float32]Portfolio)

	// HTTP request multiplexer
	// mux.Router matches incoming requests against a list of registered routes
	// and calls a handler for the route that matches the URL
	r := mux.NewRouter()

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")

	stockDealer := new(StockDealer)
	s.RegisterService(stockDealer, "")

	// middle ware: organizing  hared functionalities
	chain := alice.New(
		func(h http.Handler) http.Handler {
			return handlers.CombinedLoggingHandler(os.Stdout, h)
		},
		handlers.CompressHandler,
		func(h http.Handler) http.Handler {
			return recovery.Handler(os.Stderr, h, true)
		})

    r.Handle("/rpc", chain.Then(s))
	fmt.Println("Server listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
