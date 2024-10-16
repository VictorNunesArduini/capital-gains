package main

import (
	"bufio"
	"os"
	"strconv"
	"capital-gains/internal"
	"capital-gains/internal/application"
	"capital-gains/internal/model"
)

var propertie model.Propertie

func init()	{

	propertie = model.Propertie{
		ProfitPercentage: 0.20,
		MaxSellOperationValue: 20000.0,	
	}

	if envValue := os.Getenv(internal.ProfitPercentageEnv); envValue != "" {
		profitPercentage, _ := strconv.ParseFloat(envValue, 64)
		propertie.ProfitPercentage = profitPercentage
	}

	if envValue := os.Getenv(internal.MaxSellOperationValueEnv); envValue != "" {
		maxSellOperationValue, _ := strconv.ParseFloat(envValue, 64)
		propertie.MaxSellOperationValue = maxSellOperationValue
	}
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	line := internal.ReadStdin(reader)

	for line != "" {

		operationsInput := internal.ParseJson(line)

		taxes := application.ComputeOperations(operationsInput, propertie)

		internal.WriteStdout(taxes)

		line = internal.ReadStdin(reader)
	}
}
