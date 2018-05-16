package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/samoslab/priceservice/service"
)

//{"ok":{0|1},"data":{"price_usd":"xxx","price_btc":"xxxx"}}

// Error500 respond with a 500 error and include a message
var (
	SamosName = "samos"
)

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
	OK   int                  `json:"ok"`
	Data map[string]PriceData `json:"data"`
}

func NewPriceResult() *PriceResult {
	return &PriceResult{
		OK:   0,
		Data: map[string]PriceData{},
	}
}

type PriceData struct {
	Name     string `json:"name"`
	PriceUsd string `json:"price_usd"`
	PriceBtc string `json:"price_btc"`
	PriceCny string `json:"price_cny"`
}

type PriceManager struct {
	PriceMap map[string]PriceData
	Mutex    sync.Mutex
}

func NewPriceManager(coinTypes []string) *PriceManager {
	pm := &PriceManager{
		PriceMap: map[string]PriceData{},
	}
	for _, ct := range coinTypes {
		pm.PriceMap[ct] = PriceData{}
	}
	return pm
}

type PriceService struct {
	pm *PriceManager
}

func (ps *PriceService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	coinType := r.FormValue("coinType")
	result := NewPriceResult()
	switch coinType {
	case "bitcoin":
	case "skycoin":
	case SamosName:
	case "all":
	default:
		fmt.Printf("unsupported coin type %s\n", coinType)
		SendJSONOr500(w, result)
		return
	}
	fmt.Printf("coin type is %s\n", coinType)
	result.OK = 1
	if price, ok := ps.pm.PriceMap[coinType]; ok {
		result.Data[coinType] = price
		SendJSONOr500(w, result)
		return
	}

	for coinType, price := range ps.pm.PriceMap {
		result.Data[coinType] = price
	}
	SendJSONOr500(w, result)
}

func ConstructResponse(rsp *service.CoinMarketInfo) PriceData {
	return PriceData{
		Name:     rsp.Name,
		PriceBtc: rsp.PriceBtc,
		PriceUsd: rsp.PriceUsd,
		PriceCny: rsp.PriceCny,
	}
}

func CalcuSamosPrice(pd PriceData) PriceData {
	samosPriceData := PriceData{
		Name:     SamosName,
		PriceBtc: "0.000021",
		PriceUsd: "unknown",
		PriceCny: "unknown",
	}
	usd, err := strconv.ParseFloat(pd.PriceUsd, 10)
	if err != nil {
		return samosPriceData
	}
	cny, err := strconv.ParseFloat(pd.PriceCny, 10)
	if err != nil {
		return samosPriceData
	}
	samosUsd := usd * 0.000021
	samosPriceData.PriceUsd = fmt.Sprintf("%0.4f", samosUsd)
	samosCny := cny * 0.000021
	samosPriceData.PriceCny = fmt.Sprintf("%0.4f", samosCny)
	return samosPriceData
}

func CacheCoinInfo(pm *PriceManager) {
	for {
		for coinType, _ := range pm.PriceMap {
			if coinType == SamosName {
				continue
			}
			rsp, err := service.GetCoinPriceInfo(coinType)
			if err != nil {
				fmt.Printf("coinType %s get from coin market err, %v\n", coinType, err)
				continue
			}

			priceInfo := ConstructResponse(rsp)
			pm.Mutex.Lock()
			pm.PriceMap[coinType] = priceInfo
			// only for samos
			if coinType == "bitcoin" {
				pm.PriceMap[SamosName] = CalcuSamosPrice(priceInfo)
			}
			pm.Mutex.Unlock()
		}

		time.Sleep(5 * time.Second)
	}
}

func main() {

	coinTypes := []string{"bitcoin", "skycoin"}
	pm := NewPriceManager(coinTypes)

	// get coin info from coinmarket
	go CacheCoinInfo(pm)

	priceService := &PriceService{pm: pm}

	http.Handle("/price", priceService)

	if err := http.ListenAndServe(":8081", http.DefaultServeMux); err != nil {
		log.Fatalln(err)
	}
}
