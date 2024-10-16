package application

import (
	"capital-gains/internal/model"
	"capital-gains/internal"
	"capital-gains/internal/service"
	"fmt"
)

func ComputeOperations(operations []model.OperationIO, propertie model.Propertie) []model.TaxIO {

	taxes := make([]model.TaxIO, 0, len(operations))

	wallet := &model.Wallet{}

	for _, operation := range operations {

		operationAdapter := service.NewOperation(operation, wallet, propertie)

		taxValue := 0.0

		switch operation.Action {
			case model.Buy:
				operationAdapter.Buy()

			case model.Sell:
				taxValue = operationAdapter.Sell()
			
			default:
				internal.RaiseError(fmt.Sprintf("unexpected operation: %s", operation.Action))
		}

		taxes = append(taxes, model.TaxIO{taxValue})
	}

	return taxes
}
