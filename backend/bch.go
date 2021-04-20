package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const coinMarketKey = "4fde68db-6ef9-4a9f-97ed-570bb22dd3eb"

func HandleBchPrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	//var p price = "400.00"
	w.Header().Add("Access-Control-Allow-Origin", "*")

	p := GetPriceFromCoinMarket()
	out, _ := json.Marshal(p)
	_, _ = fmt.Fprintf(w, string(out))
}

type CoinQuoteResult struct {
	Data struct {
		BCH struct {
			Quote struct {
				USD struct {
					Price float64
				}
			}
		}
	}
}

func GetPriceFromCoinMarket() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest", nil)
	if err != nil {
		return ""
	}
	q := url.Values{}
	q.Add("convert", "USD")
	q.Add("symbol", "BCH")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", coinMarketKey)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to coinMarket server")
		return ""
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	c := CoinQuoteResult{}
	err = json.Unmarshal(respBody, &c)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%f", c.Data.BCH.Quote.USD.Price)
}
