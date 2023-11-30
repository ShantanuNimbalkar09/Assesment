package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"net/http"

	"github.com/gorilla/mux"
)

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
	body, err := io.ReadAll(response.Body)
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
		tickermap, ok := ticker.(map[string]interface{})
		if !ok {
			continue
		}
		symbol, ok := tickermap["ticker"].(string)
		if !ok {
			continue
		}
		stocks = append(stocks, Stock{Symbol: symbol})

	}
	for i := range stocks {
		stocks[i].RefreshInterval = rand.Intn(5) + 1
	}
	stockJSON, err := json.Marshal(stocks)
	if err != nil {
		http.Error(w, "Failed to marshal stocks to JSON", http.StatusInternalServerError)
		return
	}
	err = os.WriteFile("stocks.json", stockJSON, 0644)
	if err != nil {
		http.Error(w, "Failed to write stocks to file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
	// fmt.Fprintln(w, "Stocks fetched successfully!")
}

func fetchPreviousClose(w http.ResponseWriter, r *http.Request) {
	stocksJSON, err := os.ReadFile("stocks.Json")
	if err != nil {
		http.Error(w, "Failed to read stocks from file", http.StatusInternalServerError)
		return
	}
	var stocks []Stock
	if err := json.Unmarshal(stocksJSON, &stocks); err != nil {
		http.Error(w, "Failed to unmarshal stocks from JSON", http.StatusInternalServerError)
		return
	}

	var prices []Price
	for _, stock := range stocks {
		prices = append(prices, Price{Name: stock.Symbol, Price: rand.Float64() * 100})
	}

	PriceJSON, err := json.Marshal(prices)
	if err != nil {
		http.Error(w, "Failed to marshal prices to JSON", http.StatusInternalServerError)
		return
	}
	err = os.WriteFile("prices.json", PriceJSON, 0644)
	if err != nil {
		http.Error(w, "Failed to write prices to file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prices)
}

func updatePrices() {

	for {

		pricesJSON, err := os.ReadFile("Prices.json")
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

		for i := range prices {
			prices[i].Price += rand.Float64()*5 - 2.5
		}
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

	router := mux.NewRouter()

	// Register the "/fetch-stocks" endpoint
	router.HandleFunc("/fetch-stocks/{n:[0-9]+}", fetchStocks).Methods("GET")
	router.HandleFunc("/fetch-previous-close", fetchPreviousClose).Methods("GET")

	go updatePrices()
	// Define the port to listen on
	port := 3000

	// Print a message indicating that the server is running
	fmt.Printf("Server is running at http://localhost:%d\n", port)

	// Start the HTTP server
	http.Handle("/", router)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
