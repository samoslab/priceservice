package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/samoslab/priceservice/service"
)

//{"ok":{0|1},"data":{"price_usd":"xxx","price_btc":"xxxx"}}

// Error500 respond with a 500 error and include a message
func Error500(w http.ResponseWriter, msg string) {
	errorXXXMsg(w, http.StatusInternalServerError, msg)
}

func errorXXXMsg(w http.ResponseWriter, status int, msg string) {
	httpMsg := http.StatusText(status)
	if msg != "" {
		httpMsg = fmt.Sprintf("%s - %s", httpMsg, msg)
	}
	HTTPError(w, status, httpMsg)
}

// HTTPError wraps http.Error
func HTTPError(w http.ResponseWriter, status int, httpMsg string) {
	msg := fmt.Sprintf("%d %s", status, httpMsg)
	http.Error(w, msg, status)
}

// SendJSONOr500 writes an object as JSON, writing a 500 error if it fails
func SendJSONOr500(w http.ResponseWriter, m interface{}) {
	out, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		Error500(w, "json.MarshalIndent failed")
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if _, err := w.Write(out); err != nil {
		fmt.Printf("http Write failed %v\n", err)
	}
}

type PriceResult struct {
	OK   int `json:"ok"`
	Data struct {
		PriceUsd string `json:"price_usd"`
		PriceBtc string `json:"price_btc"`
		PriceCny string `json:"price_cny"`
	} `json:"data"`
}

type PriceManager struct {
	PriceMap map[string]PriceResult
	Mutex    sync.Mutex
}

func HandlePrice(w http.ResponseWriter, r *http.Request) {
	coinType := r.FormValue("coinType")
	rsp := &service.CoinMarketInfo{}
	result := PriceResult{}
	result.OK = 0
	var err error
	switch coinType {
	case "bitcoin":
		rsp, err = service.GetCoinPriceInfo("bitcoin")
	case "skycoin":
		rsp, err = service.GetCoinPriceInfo("skycoin")
	default:
	}
	if err != nil {
		fmt.Fprintf(w, "coinType %s err, %v", coinType, err)
		return
	}

	result.OK = 1
	result.Data.PriceBtc = rsp.PriceBtc
	result.Data.PriceUsd = rsp.PriceUsd
	result.Data.PriceCny = rsp.PriceCny
	SendJSONOr500(w, result)
}

func main() {

	http.HandleFunc("/price", HandlePrice)

	if err := http.ListenAndServe(":8081", http.DefaultServeMux); err != nil {
		log.Fatalln(err)
	}
}
