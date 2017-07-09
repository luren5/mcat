package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/luren5/mcat/ethabi"
	"github.com/luren5/mcat/utils"
	"github.com/spf13/cobra"
)

var (
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
		fmt.Println("func def:", funcDef)
		fmt.Println("func selector:", selector)
		fmt.Println("call bytes:", callBytes)
		// params bytes code
	},
}

func init() {
	RootCmd.AddCommand(callCmd)

	callCmd.Flags().StringVar(&contract, "contract", "", "The contract name.")
	callCmd.Flags().StringVar(&function, "function", "", "the function name.")
	callCmd.Flags().StringVar(&params, "params", "", "params")
}
