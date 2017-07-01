package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/luren5/mcat/common"
)

var compiledDir string = "compiled"

func CompiledDir() string {
	pwd, _ := os.Getwd()
	return pwd + "/" + compiledDir + "/"
}

func ConfigPath() string {
	pwd, _ := os.Getwd()
	return pwd + "/mcat.yaml"
}

func CheckIfTxMined(ip, rpc_port, txHash string) (interface{}, error) {
	params := fmt.Sprintf(`"%s"`, txHash)
	res, err := JrpcPost(ip, rpc_port, "eth_getTransactionReceipt", params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SendTransaction(ip, rpc_port string, tx *common.Transaction) (interface{}, error) {
	var params string
	switch tx.Type {
	case common.TxTypeCommon:
		params = fmt.Sprintf(`{"from": "%s", "to": "%s", "gas": "%s", "gasPrice": "%s","value": "%s", "data": "%s"}`, tx.From, tx.To, tx.Gas, tx.GasPrice, tx.Value, tx.Data)
	case common.TxTypeContract:
		params = fmt.Sprintf(`{"from": "%s", "gas": "%s", "gasPrice": "%s","value": "%s", "data": "%s"}`, tx.From, tx.Gas, tx.GasPrice, tx.Value, tx.Data)
	default:
		return nil, errors.New("Invalid tx type")

	}
	fmt.Println(params)
	return JrpcPost(ip, rpc_port, "eth_sendTransaction", params)
}
