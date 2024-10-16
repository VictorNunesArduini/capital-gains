package service

import (
	"capital-gains/internal/model"
	"capital-gains/internal"
)

type OperationAdapter interface {
	Buy()
	Sell() float64
}

type OperationModel struct {
	Action string
	StockValue float64
	Quantity int32
	Wallet *model.Wallet
	Propertie model.Propertie
}

func NewOperation(op model.OperationIO, wallet *model.Wallet, propertie model.Propertie) OperationAdapter {
	return &OperationModel {
		Action: op.Action,
		StockValue: op.StockValue,
		Quantity: op.Quantity,
		Wallet: wallet,
		Propertie: propertie,
	}
}

func (op *OperationModel) Buy() {

	if op.Wallet.AccumulatedQuantity == 0.0 {
		op.Wallet.Difference = 0.0
	}

	op.Wallet.AverageStockValue = ((float64(op.Wallet.AccumulatedQuantity) * op.Wallet.AverageStockValue) + (float64(op.Quantity) * op.StockValue)) / (float64(op.Wallet.AccumulatedQuantity + op.Quantity))
	//500 * 20.0 + 10.0 * 1000.0 / 1500 = 
	op.Wallet.AccumulatedQuantity = op.Wallet.AccumulatedQuantity + op.Quantity
}

func (op *OperationModel) Sell() float64 {

	if op.Quantity > op.Wallet.AccumulatedQuantity {
		internal.RaiseError("quantity to sell is greater than accumulated quantity!")
	}

	op.Wallet.AccumulatedQuantity = op.Wallet.AccumulatedQuantity - op.Quantity

	return calculateTax(op)
}

func calculateTax(op *OperationModel) float64 {

	sellOperationValue := float64(op.Quantity) * op.StockValue 

	difference := sellOperationValue - (float64(op.Quantity) * op.Wallet.AverageStockValue)

	op.Wallet.Difference = op.Wallet.Difference + difference

	if op.Wallet.Difference <= 0.0 || sellOperationValue <= op.Propertie.MaxSellOperationValue {
		return 0.0
	}

	profit := op.Propertie.ProfitPercentage * op.Wallet.Difference

	return internal.RoundToTwoDecimals(profit)
}