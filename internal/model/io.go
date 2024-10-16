package model

type OperationIO struct {
	Action string `json:"operation"`
	StockValue float64 `json:"unit-cost"`
	Quantity int32 `json:"quantity"`
}

type TaxIO struct {
	Value float64 `json:"tax"`
}