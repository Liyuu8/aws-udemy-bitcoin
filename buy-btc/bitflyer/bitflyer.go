package bitflyer

import (
	"buy-btc/utils"
	"encoding/hex"
	"encoding/json"
	"errors"

	"crypto/hmac"
	"crypto/sha256"
	"strconv"
	"time"
)

const baseURL = "https://api.bitflyer.com"
const productCodeKey = "product_code"
const btcMinimumAmount = 0.001 // bitflyerのBTC最小注文数量
const btcMinimumAmountPlace = 4.0

type APIClient struct {
	apiKey    string
	apiSecret string
}

func NewAPIClient(apiKey, apiSecret string) *APIClient {
	return &APIClient{apiKey, apiSecret}
}

// 価格取得機能
func GetTicker(tickerChan chan *Ticker, errChan chan error, code ProductCode) {
	url := baseURL + "/v1/ticker"
	res, err := utils.DoHttpRequest("GET", url, nil, map[string]string{productCodeKey: code.String()}, nil)
	if err != nil {
		tickerChan <- nil
		errChan <- err
		return
	}

	var ticker Ticker
	err = json.Unmarshal(res, &ticker)
	if err != nil {
		tickerChan <- nil
		errChan <- err
		return
	}

	tickerChan <- &ticker
	errChan <- nil
}

func (client *APIClient) getHeader(method, path string, body []byte) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// ACCESS-SIGN は、ACCESS-TIMESTAMP, HTTP メソッド, リクエストのパス, リクエストボディを
	// 文字列として連結したものを、
	sign := timestamp + method + path + string(body)
	// API secret で HMAC-SHA256 署名を行った結果です。
	mac := hmac.New(sha256.New, []byte(client.apiSecret))
	mac.Write([]byte(sign))
	encodedSign := hex.EncodeToString(mac.Sum(nil))

	return map[string]string{
		"ACCESS-KEY":       client.apiKey,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      encodedSign,
		"Content-Type":     "application/json",
	}
}

// 新規注文作成
func (client *APIClient) makeOrderResponse(order *Order) (*OrderRes, error) {
	method := "POST"
	path := "/v1/me/sendchildorder"
	url := baseURL + path
	data, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	header := client.getHeader(method, path, data)

	res, err := utils.DoHttpRequest(method, url, header, map[string]string{}, data)
	if err != nil {
		return nil, err
	}

	var orderRes OrderRes
	err = json.Unmarshal(res, &orderRes)
	if err != nil {
		return nil, err
	}

	if len(orderRes.ChildOrderAcceptanceId) == 0 {
		return nil, errors.New(string(res))
	}

	return &orderRes, nil
}

// 新規注文実行
func (client *APIClient) PlaceOrder(price, size float64) (*OrderRes, error) {
	order := Order{
		ProductCode:     BtcJpy.String(),
		ChildOrderType:  Limit.String(),
		Side:            Buy.String(),
		Price:           price,
		Size:            size,
		MinuteToExpires: 1440, // 1day
		TimeInForce:     Gtc.String(),
	}

	orderRes, err := client.makeOrderResponse(&order)
	if err != nil {
		return nil, err
	}

	return orderRes, nil
}

// ロジックを取得する関数
func GetBuyLogic(strategy int) func(float64, *Ticker) (float64, float64) {
	var logic func(float64, *Ticker) (float64, float64)

	// TODO: Add strategies
	switch strategy {
	case 1:
		// LTP の 98.5% の価格で購入
		logic = func(budget float64, ticker *Ticker) (float64, float64) {
			buyPrice := utils.RoundDecimal(ticker.Ltp * 0.985)
			buySize := utils.CalcAmount(buyPrice, budget, btcMinimumAmount, btcMinimumAmountPlace)

			// TODO: Add currency pairs
			// if ticker.ProductCode == BtcJpy.String() {
			// 	xxx
			// } else if ticker.ProductCode == EthJpy.String() {
			// 	xxx
			// }
			return buyPrice, buySize
		}
	default:
		// BestASK で購入
		logic = func(budget float64, ticker *Ticker) (float64, float64) {
			buyPrice := utils.RoundDecimal(ticker.BestAsk)
			buySize := utils.CalcAmount(buyPrice, budget, btcMinimumAmount, btcMinimumAmountPlace)

			return buyPrice, buySize
		}
	}

	return logic
}

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

type Order struct {
	ProductCode     string  `json:"product_code"`
	ChildOrderType  string  `json:"child_order_type"`
	Side            string  `json:"side"`
	Price           float64 `json:"price"`
	Size            float64 `json:"size"`
	MinuteToExpires int     `json:"minute_to_expire"`
	TimeInForce     string  `json:"time_in_force"`
}

type OrderRes struct {
	ChildOrderAcceptanceId string `json:"child_order_acceptance_id"`
}
