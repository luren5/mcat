package cmd

import (
	"fmt"
	"os"

	"github.com/luren5/mcat/common"
	"github.com/luren5/mcat/utils"
	"github.com/spf13/cobra"
)

// gasPriceCmd represents the gasPrice command
var gasPriceCmd = &cobra.Command{
	Use:   "gasPrice",
	Short: "Show the current gas price.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ip, rpc_port, err := utils.GetRpcInfo()
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		gasPrice, err := common.GasPrice(ip, rpc_port)
		if err != nil {
			fmt.Printf("Failed to get gas price, %v \r\n", err)
			os.Exit(-1)
		}
		fmt.Printf("The current gas price is %s \r\n", gasPrice)
	},
}

func init() {
	RootCmd.AddCommand(gasPriceCmd)
}
