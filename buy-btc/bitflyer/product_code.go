package bitflyer

type ProductCode int

const (
	BtcJpy ProductCode = iota
	EthJpy
	FxBtcJpy
	EthBtc
	BchBtc
)

func (code ProductCode) String() string {
	switch code {
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
