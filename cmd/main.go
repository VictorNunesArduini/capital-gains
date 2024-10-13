package main

import (
	"fmt"
    "strings"
	"encoding/json"
	"io"
	"math"
	"os"
	"bufio"
	"runtime"
	"strconv"
)

//
//nova-média-ponderada = ((quantidade-de-ações-atual * média-ponderada-atual) + (quantidade-de-ações-compradas * valor-de-compra)) / (quantidade-de-ações-atual + quantidade-de-ações-compradas)
//

// desafios encontrados: 
// como lidar com floats para valores financeiros?
// utilizar env vars - ok
// pensar em uma estrutura que suporte mais de uma ação e manter o stockWallet para ela - ok
// json identation - ok


type Operation struct {
	Stock string `json:"stock, omitempty"`
	Action string `json:"operation"`
	StockValue float64 `json:"unit-cost"`
	Quantity int32 `json:"quantity"`
}

type Tax struct {
	Value float64 `json:"tax"`
}

type stockWallet struct {
	AverageStockValue float64
 	AccumulatedQuantity int32
	Difference float64
}

var stock map[string]*stockWallet

var profitPercentage float64 = 0.20
var maxSellOperationValue float64 = 20000.0


func init()	{

	if envValue := os.Getenv("PROFIT_PERCENTAGE"); envValue != "" {
		profitPercentage, _ = strconv.ParseFloat(envValue, 64)
		fmt.Printf("profit percentage: %f\n", profitPercentage)
	}

	if envValue := os.Getenv("MAX_SELL_OPERATION_VALUE"); envValue != "" {
		maxSellOperationValue, _ = strconv.ParseFloat(envValue, 64)
	}
}

func main()  {

	reader := bufio.NewReader(os.Stdin)

	line := readStdin(reader)

	for line != "" {

		stock := make(map[string]*stockWallet)  // reset the stocks wallet

		operations := parseJson(line)
		taxes := make([]Tax, 0, len(operations))

		for _, operation := range operations {

			var taxValue float64 = 0.0

			if operation.Action == "buy" {

				if _, ok := stock[operation.Stock]; !ok {
					stock[operation.Stock] = &stockWallet{}
				}

				operation.buy(stock[operation.Stock])
			} else if operation.Action == "sell" {
				
				if _, ok := stock[operation.Stock]; !ok {
					fmt.Errorf("we can't sell this inexistent stock: %s", operation.Stock)
					os.Exit(1)
				}

				taxValue = operation.sell(stock[operation.Stock])
			}

			taxes = append(taxes, Tax{taxValue})
		}

		writeStdout(taxes)

		line = readStdin(reader)
	}
}

func readStdin(reader *bufio.Reader) string {

	out, _, err := reader.ReadLine()
    if err == io.EOF {
        return ""
    }

	return strings.TrimRight(string(out), newLineByOs())
}

func newLineByOs() string {

	if runtime.GOOS == "windows" {
		return "\r\n"
	} else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		return "\n"
	}

	return "\r"
}

func writeStdout(taxes []Tax) {

	fmt.Println(stringfyJson(taxes))
}

func parseJson(line string) []Operation {

	var operations []Operation

	err := json.Unmarshal([]byte(line), &operations)

	if err != nil {
        return []Operation{}
    }

	return operations
}

func stringfyJson(taxes []Tax) string {
	jsonTaxes, err := json.MarshalIndent(taxes, "", "")

	if err != nil {
		fmt.Errorf("failed to stringfy json: %v", taxes)
		os.Exit(1)
	}

	return strings.ReplaceAll(string(jsonTaxes), newLineByOs(), "")
}

func (op Operation) buy(wallet *stockWallet) {

	wallet.AverageStockValue = ((float64(wallet.AccumulatedQuantity) * wallet.AverageStockValue) + (float64(op.Quantity) * op.StockValue)) / (float64(wallet.AccumulatedQuantity + op.Quantity))

	wallet.AccumulatedQuantity = wallet.AccumulatedQuantity + op.Quantity
}

func (op Operation) sell(wallet *stockWallet) float64 {

	if op.Quantity > wallet.AccumulatedQuantity {
		fmt.Errorf("quantity is greater than accumulated quantity!")
		os.Exit(1)
	}

	wallet.AccumulatedQuantity = wallet.AccumulatedQuantity - op.Quantity

	return calculateTax(op, wallet)
}

func calculateTax(op Operation, wallet *stockWallet) float64 {

	sellOperationValue := float64(op.Quantity) * op.StockValue

	difference := sellOperationValue - (float64(op.Quantity) * wallet.AverageStockValue)

	wallet.Difference = wallet.Difference + difference

	if wallet.Difference <= 0.0 || sellOperationValue <= maxSellOperationValue {
		return 0.0
	}

	profit := profitPercentage * wallet.Difference

	return roundToTwoDecimals(profit)
}


func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}