package cmd

import (
	"fmt"
	"os"

	"github.com/luren5/mcat/common"
	"github.com/luren5/mcat/ethabi"
	"github.com/luren5/mcat/utils"
	"github.com/spf13/cobra"
)

var blockNum string

// callConstCmd represents the callConst command
var callConstCmd = &cobra.Command{
	Use:   "callConst",
	Short: "Call the contract function without sending a transaction.",
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

		// ethabi
		e, err := ethabi.NewEthABI(contract)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		// is constant
		if b, err := e.IsConstant(function); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		} else {
			if !b {
				fmt.Printf("%s is not a constant function")
				os.Exit(-1)
			}
		}
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
		fmt.Println(callBytes)
		// call
		ip, rpc_port, err := utils.GetRpcInfo()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		if r, err := common.EthCall(ip, rpc_port, addr, callBytes, blockNum); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(r)
		}
	},
}

func init() {
	RootCmd.AddCommand(callConstCmd)
	callConstCmd.Flags().StringVar(&blockNum, "blockNum", "latest", "the block num")
	callConstCmd.Flags().StringVar(&addr, "addr", "", "The contract address.")
	callConstCmd.Flags().StringVar(&contract, "contract", "", "The contract name.")
	callConstCmd.Flags().StringVar(&function, "function", "", "the function name.")
	callConstCmd.Flags().StringVar(&params, "params", "", "params")
}
