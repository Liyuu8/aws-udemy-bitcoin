package utils

import (
	"math"
)

func RoundDecimal(num float64) float64 {
	return math.Round(num)
}

// 少数の切り上げ機能
func roundUp(num, places float64) float64 {
	shift := math.Pow(10, places)
	return RoundDecimal(num*shift) / shift
}

// 予算をもとに購入数量を算出する機能
func CalcAmount(price, budget, minimumAmount, places float64) float64 {
	amout := roundUp(budget/price, places)
	if amout < minimumAmount {
		return minimumAmount
	} else {
		return amout
	}
}
