package bitflyer

type ProductCode int

const (
	BtcJpy ProductCode = iota
	EthJpy
	FxBtcJpy
	EthBtc
	BchBtc
)

type OrderType int

// 指値 or 成行
const (
	Limit OrderType = iota
	Market
)

type Side int

// 買い or 売り
const (
	Buy Side = iota
	Sell
)

type TimeInForce int

// 執行数量条件
const (
	Gtc TimeInForce = iota
	Ioc
	Fok
)

func (productCode ProductCode) String() string {
	switch productCode {
	case BtcJpy:
		return "BTC_JPY"
	case EthJpy:
		return "ETH_JPY"
	case FxBtcJpy:
		return "FX_BTC_JPY"
	case EthBtc:
		return "ETH_BTC"
	case BchBtc:
		return "BCH_BTC"
	default:
		return "BTC_JPY"
	}
}

func (orderType OrderType) String() string {
	switch orderType {
	case Limit:
		return "LIMIT"
	case Market:
		return "MARKET"
	default:
		return "LIMIT"
	}
}

func (side Side) String() string {
	switch side {
	case Buy:
		return "BUY"
	case Sell:
		return "SELL"
	default:
		return "BUY"
	}
}

func (timeInForce TimeInForce) String() string {
	switch timeInForce {
	case Gtc:
		return "GTC"
	case Ioc:
		return "IOC"
	case Fok:
		return "FOK"
	default:
		return "GTC"
	}
}
