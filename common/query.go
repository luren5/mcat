package common

import (
	"fmt"

	"github.com/luren5/mcat/utils"
)

func EthCall(ip, rpc_port, contractAddr, callBytes, blockNum string) (string, error) {
	params := fmt.Sprintf(`{"to":"%s","data":"%s"},"%s"`, contractAddr, callBytes, blockNum)
	res, err := utils.JrpcPost(ip, rpc_port, "eth_call", params)
	if err != nil {
		return "", err
	}
	return res.(string), err
}
