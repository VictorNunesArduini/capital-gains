package internal

import (
	"math"
	"fmt"
	"os"
)

func RoundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

func RaiseError(message string) {
	fmt.Errorf(message)
	os.Exit(1)
}