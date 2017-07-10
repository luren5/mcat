package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luren5/mcat/common"
	"github.com/luren5/mcat/ethabi"
	"github.com/luren5/mcat/utils"
	"github.com/spf13/cobra"
)

var (
	addr     string
	contract string
	function string
	params   string
)

// callCmd represents the call command
var callCmd = &cobra.Command{
	Use:   "call",
	Short: "Call contract function.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(contract) == 0 {
			fmt.Println("Invalid contract name")
			os.Exit(-1)
		}
		if len(function) == 0 {
			fmt.Println("Invalid function name")
			os.Exit(-1)
		}

		abiFile := utils.CompiledDir() + contract + ".abi"
		abiBytes, err := ioutil.ReadFile(abiFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		e := ethabi.NewEthABI(contract, abiBytes)

		// selector
		funcDef, err := e.FuncDef(function)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		selector := ethabi.CalSelector(funcDef)

		// params
		cp, err := e.CombineParams(function, params)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		callBytes, err := ethabi.CalBytes(selector, cp)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		// send tx
		ip, rpc_port, err := utils.GetRpcInfo()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		account, err := utils.GetDefaultAccount()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		tx := new(common.Transaction)
		tx.From = account
		tx.To = addr
		tx.Data = callBytes
		tx.Type = common.TxTypeCommon

		// gas
		if tx.GasPrice, err = common.GasPrice(ip, rpc_port, tx); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		if tx.Gas, err = common.EstimateGas(ip, rpc_port, tx); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		hash, err := common.SendTransaction(ip, rpc_port, tx)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Succeed in calling %s, tx hash: %s \r\n", function, hash)
	},
}

func init() {
	RootCmd.AddCommand(callCmd)

	callCmd.Flags().StringVar(&addr, "addr", "", "The contract address.")
	callCmd.Flags().StringVar(&contract, "contract", "", "The contract name.")
	callCmd.Flags().StringVar(&function, "function", "", "the function name.")
	callCmd.Flags().StringVar(&params, "params", "", "params")
}
