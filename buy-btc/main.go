package main

import (
	"buy-btc/bitflyer"
	"fmt"
	"math"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	apiKey, err := getParameter("buy_btc_api_key")
	if err != nil {
		return getErrorResponse(err.Error()), err
	}

	apiSecret, err := getParameter("buy_btc_api_secret")
	if err != nil {
		return getErrorResponse(err.Error()), err
	}

	client := bitflyer.NewAPIClient(apiKey, apiSecret)

	tickerChan := make(chan *bitflyer.Ticker)
	errChan := make(chan error)
	defer close(tickerChan)
	defer close(errChan)

	go bitflyer.GetTicker(tickerChan, errChan, bitflyer.BtcJpy)
	ticker := <-tickerChan
	err = <-errChan
	if err != nil {
		return getErrorResponse(err.Error()), err
	}

	order := bitflyer.Order{
		ProductCode:     bitflyer.BtcJpy.String(),
		ChildOrderType:  bitflyer.Limit.String(),
		Side:            bitflyer.Buy.String(),
		Price:           math.Round(ticker.Ltp * 0.90), // 現在価格の90%
		Size:            0.001,                         // 0.001BTC
		MinuteToExpires: 1440,                          // 1day
		TimeInForce:     bitflyer.Gtc.String(),
	}

	orderRes, err := client.PlaceOrder(&order)
	if err != nil {
		return getErrorResponse(err.Error()), err
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Order Response: %+v", orderRes),
		StatusCode: 200,
	}, nil
}

// System Manager からパラメータを取得する関数
func getParameter(key string) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		// ローカル環境設定を読み込む (~/.aws/config)
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ssm.New(sess, aws.NewConfig().WithRegion("ap-northeast-1"))

	params := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}

	res, err := svc.GetParameter(params)
	if err != nil {
		return "", err
	}

	return *res.Parameter.Value, nil
}

func getErrorResponse(message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: 400,
	}
}

func main() {
	lambda.Start(handler)
}
