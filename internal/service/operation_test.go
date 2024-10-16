package service

import (
	"testing"
	"capital-gains/internal/model"
	"capital-gains/internal"
)

func TestBuy(t *testing.T) {

	defaultPropertie := model.Propertie{
		ProfitPercentage: 0.20,
		MaxSellOperationValue: 20000.0,
	}

	buyOp := model.OperationIO{
		Action: "buy",
		StockValue: 10.00,
		Quantity: 1000,
	}

	wallet := &model.Wallet{
		AverageStockValue: 20.0,
		AccumulatedQuantity: 500,
		Difference: 0.0,
	}

	operationAdapter := NewOperation(buyOp, wallet, defaultPropertie)

	//unit test
    operationAdapter.Buy()

    // Assert expected results
    expectedAverageStockValue := 13.33
    var expectedAccumulatedQuantity int32 = 1500

    if internal.RoundToTwoDecimals(wallet.AverageStockValue) != expectedAverageStockValue {
        t.Errorf("Expected AverageStockValue to be %f, got %f", expectedAverageStockValue, wallet.AverageStockValue)
    }

    if wallet.AccumulatedQuantity != expectedAccumulatedQuantity {
        t.Errorf("Expected AccumulatedQuantity to be %d, got %d", expectedAccumulatedQuantity, wallet.AccumulatedQuantity)
    }
}


//simple test, we could build a structure that simulates a variety of scenarios
func TestSell(t *testing.T) {

	defaultPropertie := model.Propertie{
		ProfitPercentage: 0.20,
		MaxSellOperationValue: 20000.0,
	}

    sellOp := model.OperationIO{
		Action: "sell",
		StockValue: 50.00,
		Quantity: 1000,
	}

	wallet := &model.Wallet{
		AverageStockValue: 13.33,
		AccumulatedQuantity: 1500,
		Difference: 0.0,
	}

	operationAdapter := NewOperation(sellOp, wallet, defaultPropertie)

	//unit test
	taxValue := operationAdapter.Sell()

	// Assert expected results
	expectedDifference := 36670.00
	expectedTaxValue := 7334.00
	var expectedAccumulatedQuantity int32 = 500

	if wallet.AccumulatedQuantity != expectedAccumulatedQuantity {
        t.Errorf("Expected AccumulatedQuantity to be %d, got %d", expectedAccumulatedQuantity, wallet.AccumulatedQuantity)
    }

	if internal.RoundToTwoDecimals(taxValue) != expectedTaxValue {
        t.Errorf("Expected taxValue to be %f, got %f", expectedTaxValue, taxValue)
    }

	if internal.RoundToTwoDecimals(wallet.Difference) != expectedDifference {
        t.Errorf("Expected Difference to be %f, got %f", expectedDifference, wallet.Difference)
    }
}
