package common

import (
	"errors"
	"fmt"

	"github.com/luren5/mcat/utils"
)

type Transaction struct {
	From     string
	To       string
	Gas      string
	GasPrice string
	Value    string
	Data     string
	Type     uint
}

const (
	TxTypeCommon = iota
	TxTypeDeploy
)

func CheckIfTxMined(ip, rpc_port, txHash string) (interface{}, error) {
	params := fmt.Sprintf(`"%s"`, txHash)
	res, err := utils.JrpcPost(ip, rpc_port, "eth_getTransactionReceipt", params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func EstimateGas(ip, rpc_port string, tx *Transaction) (string, error) {
	param, err := generateTxParam(tx)
	if err != nil {
		return "", err
	}
	r, err := utils.JrpcPost(ip, rpc_port, "eth_estimateGas", param)
	if err != nil {
		return "", err
	}
	return r.(string), err
}

func generateTxParam(tx *Transaction) (string, error) {
	var params string
	switch tx.Type {
	case TxTypeCommon:
		params = fmt.Sprintf(`{"from": "%s", "to": "%s", "gas": "%s", "gasPrice": "%s","value": "%s", "data": "%s"}`, tx.From, tx.To, tx.Gas, tx.GasPrice, tx.Value, tx.Data)
	case TxTypeDeploy:
		params = fmt.Sprintf(`{"from": "%s", "gas": "%s", "gasPrice": "%s","value": "%s", "data": "%s"}`, tx.From, tx.Gas, tx.GasPrice, tx.Value, tx.Data)
		//fmt.Println(params)
	default:
		return "", errors.New("Invalid tx type")

	}
	return params, nil
}

func SendTransaction(ip, rpc_port string, tx *Transaction) (interface{}, error) {
	param, err := generateTxParam(tx)
	if err != nil {
		return nil, err
	}
	return utils.JrpcPost(ip, rpc_port, "eth_sendTransaction", param)
}

func GasPrice(ip, rpc_port string, tx *Transaction) (string, error) {
	r, err := utils.JrpcPost(ip, rpc_port, "eth_gasPrice", "")
	if err != nil {
		return "", err
	}
	return r.(string), err
}
