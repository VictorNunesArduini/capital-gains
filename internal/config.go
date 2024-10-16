package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strings"
	"capital-gains/internal/model"
)

const (
	windowsOs = "windows"
	linuxOs = "linux"
	macOs = "darwin"
)

func ReadStdin(reader *bufio.Reader) string {

	out, _, err := reader.ReadLine()
	
    if err == io.EOF {
        return ""
    }

	return strings.TrimRight(string(out), newLineByOs())
}

func WriteStdout(taxes []model.TaxIO) {

	fmt.Println(stringfyJson(taxes))
}

func ParseJson(input string) []model.OperationIO {

	var operations []model.OperationIO

	err := json.Unmarshal([]byte(input), &operations)

	if err != nil {
        return []model.OperationIO{}
    }

	return operations
}

func stringfyJson(taxes []model.TaxIO) string {
	jsonTaxes, err := json.MarshalIndent(taxes, "", "")

	if err != nil {
		RaiseError(fmt.Sprintf("failed to stringfy json: %v", taxes))
	}

	return strings.ReplaceAll(string(jsonTaxes), newLineByOs(), "")
}

func newLineByOs() string {

	if runtime.GOOS == linuxOs || runtime.GOOS == macOs {
		return "\n"
	} else if runtime.GOOS == windowsOs {
		return "\r\n"
	}

	return "\r"
}
