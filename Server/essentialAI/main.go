// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	polygon "github.com/polygon-io/client-go/rest"
// 	"github.com/polygon-io/client-go/rest/models"
// )

// func main() {

// 	c := polygon.New("XhTXpb9ed9p5QVhmQ3BC1rS9Hywp78e7")

// 	params := models.GetTickerDetailsParams{
// 		Ticker: "AAPL",

// 	}.WithDate(models.Date(time.Date(2021, 7, 22, 0, 0, 0, 0, time.Local)))

// 	//Ticker aggregate bars whole data
// 	res, err := c.GetTickerDetails(context.Background(), params)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	d := models.GetTickerDetailsResponse{
// 		BaseResponse: res.BaseResponse,
// 		Results:      res.Results,
// 	}

// 	fmt.Printf("%T", d)

// 	jsonData, err := json.Marshal(d.Results.Name)
// 	fmt.Println(string(jsonData))

// 	// fmt.Printf("%T", res)

// 	// //it gets the previous closed data of followed Ticker
// 	// data, err := c.HTTP.JSONMarshal(res)
// 	// fmt.Printf("%T", data)
// 	// fmt.Println(string(data))

// 	// fmt.Println("                             ------------------------ ---------------- ")

// 	// prev := models.GetPreviousCloseAggParams{
// 	// 	Ticker: "AAPL",
// 	// }.WithAdjusted(true)

// 	// close, err := c.GetPreviousCloseAgg(context.Background(), prev)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// d := models.GetPreviousCloseAggResponse{
// 	// 	Ticker: "AAPL",
// 	// }

// 	//fmt.Printf("%T", close)
// 	// prevclose := c.HTTP.JSONUnmarshal()
// 	// fmt.Printf("%T", prevclose)
// 	// fmt.Println(string(prevclose))
// 	// fmt.Println(close.NextURL)
// 	// fmt.Println(string(prevclose))
// 	// fmt.Println("                             ------------------------ ---------------- ")

// 	//finding Tickers

// 	// tickers := models.GetAllTickersSnapshotParams{

// 	// }.MarketType
// 	// details, err := c.GetAllTickersSnapshot(context.Background(), tickers)
// 	// tickerdetails, err := c.HTTP.JSONMarshal(details)
// 	// fmt.Println(string(tickerdetails))

// 	// fmt.Println("                             ------------------------ ---------------- ")
// 	// // tickernames,err:=c.GetTickerDetails()
// 	// tickernews := models.TickerNews{}
// 	// news, err := c.HTTP.JSONMarshal(tickernews)
// 	// fmt.Println(string(news))

// 	//Getting daily trade for getting tickers

// 	// dailytrade := models.GetGroupedDailyAggsParams{
// 	// 	Date: models.Date(time.Date(2021, 2, 23, 0, 0, 0, 0, time.Local)),
// 	// }.WithAdjusted(true)
// 	// dailydata, err := c.GetGroupedDailyAggs(context.Background(), dailytrade)
// 	// fmt.Println(dailydata)
// 	// s, err := c.HTTP.JSONMarshal(dailydata.Ticker)

// 	// fmt.Printf("%+v\n", dailydata)

// 	// // Print the Ticker field directly
// 	// fmt.Printf("%+v\n", dailydata.Ticker)

// 	// fmt.Println(string(s))

// 	ticker := models.Ticker{

// 	}
// 	response, err := c.GetAllTickersSnapshot(context.Background(), &ticker)
// 	s := models.GetAllTickersSnapshotResponse{
// 		Tickers: response.Tickers,

// 	}
// 	ans, err := json.Marshal(s.Tickers)
// 	fmt.Println(string(ans))
// }

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Stock represents the structure of a stock
type Stock struct {
	Symbol          string `json:"symbol"`
	RefreshInterval int    `json:"refreshInterval"`
}

// Price represents the structure of a stock price
type Price struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var polygonAPIKey = "XhTXpb9ed9p5QVhmQ3BC1rS9Hywp78e7"

func fetchStocks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n := vars["n"]

	url := fmt.Sprintf("https://api.polygon.io/v3/reference/tickers?limit=%s&apiKey=%s", n, polygonAPIKey)

	response, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch stocks from Polygon API", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Failed to parse response body", http.StatusInternalServerError)
		return
	}

	tickers, ok := result["results"].([]interface{})
	if !ok {
		http.Error(w, "Unexpected response format", http.StatusInternalServerError)
		return
	}

	var stocks []Stock
	for _, ticker := range tickers {
		tickerMap, ok := ticker.(map[string]interface{})
		if !ok {
			continue
		}
		symbol, ok := tickerMap["ticker"].(string)
		if !ok {
			continue
		}
		stocks = append(stocks, Stock{Symbol: symbol})
	}

	// Add refresh interval
	for i := range stocks {
		stocks[i].RefreshInterval = rand.Intn(5) + 1
	}

	// Write to stocks.json
	stocksJSON, err := json.Marshal(stocks)
	if err != nil {
		http.Error(w, "Failed to marshal stocks to JSON", http.StatusInternalServerError)
		return
	}
	err = ioutil.WriteFile("stocks.json", stocksJSON, 0644)
	if err != nil {
		http.Error(w, "Failed to write stocks to file", http.StatusInternalServerError)
		return
	}

	// Respond with stocks
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
}

func fetchPreviousClose(w http.ResponseWriter, r *http.Request) {
	// Read stocks from stocks.json
	stocksJSON, err := ioutil.ReadFile("stocks.json")
	if err != nil {
		http.Error(w, "Failed to read stocks from file", http.StatusInternalServerError)
		return
	}

	var stocks []Stock
	if err := json.Unmarshal(stocksJSON, &stocks); err != nil {
		http.Error(w, "Failed to unmarshal stocks from JSON", http.StatusInternalServerError)
		return
	}

	// Create random prices
	var prices []Price
	for _, stock := range stocks {
		prices = append(prices, Price{Name: stock.Symbol, Price: rand.Float64() * 100})
	}

	// Write to prices.json
	pricesJSON, err := json.Marshal(prices)
	if err != nil {
		http.Error(w, "Failed to marshal prices to JSON", http.StatusInternalServerError)
		return
	}
	err = ioutil.WriteFile("prices.json", pricesJSON, 0644)
	if err != nil {
		http.Error(w, "Failed to write prices to file", http.StatusInternalServerError)
		return
	}

	// Respond with prices
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prices)
}

func updatePrices() {
	for {
		// Read prices from prices.json
		pricesJSON, err := ioutil.ReadFile("prices.json")
		if err != nil {
			fmt.Println("Failed to read prices from file:", err)
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		var prices []Price
		if err := json.Unmarshal(pricesJSON, &prices); err != nil {
			fmt.Println("Failed to unmarshal prices from JSON:", err)
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		// Update prices
		for i := range prices {
			prices[i].Price += rand.Float64()*5 - 2.5
		}

		// Write updated prices to prices.json
		updatedPricesJSON, err := json.Marshal(prices)
		if err != nil {
			fmt.Println("Failed to marshal updated prices to JSON:", err)
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		err = ioutil.WriteFile("prices.json", updatedPricesJSON, 0644)
		if err != nil {
			fmt.Println("Failed to write updated prices to file:", err)
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/fetch-stocks/{n:[0-9]+}", fetchStocks).Methods("GET")
	r.HandleFunc("/fetch-previous-close", fetchPreviousClose).Methods("GET")

	// Start updating prices in the background
	go updatePrices()

	port := 3000
	fmt.Printf("Server is running at http://localhost:%d\n", port)
	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
