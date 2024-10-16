package model

const (
	Buy  = "buy"
	Sell = "sell"
)

type Wallet struct {
	AverageStockValue float64
 	AccumulatedQuantity int32
	Difference float64	//This field indicates Loss or Gain and will be reseted at the stock full liquidation.
}

type Propertie struct {
	ProfitPercentage float64
	MaxSellOperationValue float64
}

