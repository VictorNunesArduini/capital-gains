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

const (
	operationBuy = "buy"
	operationSell = "sell"
	windowsOs = "windows"
	linuxOs = "linux"
	macosOs = "darwin"
)


func init()	{

	if envValue := os.Getenv("PROFIT_PERCENTAGE"); envValue != "" {
		profitPercentage, _ = strconv.ParseFloat(envValue, 64)
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

			if operation.Action == operationBuy {

				if _, ok := stock[operation.Stock]; !ok {
					stock[operation.Stock] = &stockWallet{}
				}

				operation.buy(stock[operation.Stock])
			} else if operation.Action == operationSell {
				
				if _, ok := stock[operation.Stock]; !ok {
					checkError(fmt.Sprintf("we can't sell this inexistent stock: %s", operation.Stock))
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

	if runtime.GOOS == windowsOs {
		return "\r\n"
	} else if runtime.GOOS == linuxOs || runtime.GOOS == macosOs {
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
		checkError(fmt.Sprintf("failed to stringfy json: %v", taxes))
	}

	return strings.ReplaceAll(string(jsonTaxes), newLineByOs(), "")
}

func (op Operation) buy(wallet *stockWallet) {

	if wallet.AccumulatedQuantity == 0.0 { // reset the loss/profit
		wallet.Difference = 0.0
	}

	wallet.AverageStockValue = ((float64(wallet.AccumulatedQuantity) * wallet.AverageStockValue) + (float64(op.Quantity) * op.StockValue)) / (float64(wallet.AccumulatedQuantity + op.Quantity))

	wallet.AccumulatedQuantity = wallet.AccumulatedQuantity + op.Quantity
}

func (op Operation) sell(wallet *stockWallet) float64 {

	if op.Quantity > wallet.AccumulatedQuantity {
		checkError("quantity is greater than accumulated quantity!")
	}

	wallet.AccumulatedQuantity = wallet.AccumulatedQuantity - op.Quantity

	if wallet.AccumulatedQuantity == 0.0 {
		wallet.Difference = 0.0
	}

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

func checkError(message string) {
	fmt.Errorf(message)
	os.Exit(1)
}