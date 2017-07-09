package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/luren5/mcat/common"
	"github.com/luren5/mcat/utils"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy contract",
	Long:  `create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// read bin
		binFile := utils.CompiledDir() + contract + ".bin"
		data, err := ioutil.ReadFile(binFile)
		if err != nil {
			fmt.Printf("Failed to read contract file, %v", err)
			os.Exit(-1)
		}
		bin := string(data)

		// read config
		var ip, rpc_port, account, password string
		if r, err := utils.Config("ip"); err != nil {
			fmt.Sprintf("Failed to read ip config, %v", err)
			os.Exit(-1)
		} else {
			ip = r.(string)
		}
		if r, err := utils.Config("rpc_port"); err != nil {
			fmt.Sprintf("Failed to read rpc port config, %v", err)
			os.Exit(-1)
		} else {
			rpc_port = r.(string)
		}
		if r, err := utils.Config("account"); err != nil {
			fmt.Sprintf("Invalid default account, %v", err)
			os.Exit(-1)
		} else {
			account = r.(string)
		}
		if r, err := utils.Config("password"); err != nil {
			fmt.Printf("You haven't set password of default account, please make sure the default account has been unlocked.")
		} else {
			password = r.(string)
			// unlock
			params := fmt.Sprintf(`"%s", "%s"`, account, password)
			_, err := utils.JrpcPost(ip, rpc_port, "personal_unlockAccount", params)
			if err != nil {
				fmt.Printf("Failed to unlock the default account, %v ", err)
				os.Exit(-1)
			}
		}

		tx := new(common.Transaction)
		tx.From = account
		tx.Data = bin
		tx.Type = common.TxTypeContract

		// cal gas
		tx.Gas, err = common.EstimateGas(ip, rpc_port, tx)
		if err != nil {
			fmt.Printf("Failed to estimate gas, %v ", err)
			os.Exit(-1)
		}

		// gasPrice
		tx.GasPrice, err = common.GasPrice(ip, rpc_port, tx)
		if err != nil {
			fmt.Printf("Failed to estimate gas, %v ", err)
			os.Exit(-1)
		}

		// do deploy
		txHash, err := common.SendTransaction(ip, rpc_port, tx)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Printf("Succeed in deploying contract %s, tx hash: %s. Waiting for being mined…\r\n \r\n", contract, txHash)

		// check status
		for {
			time.Sleep(time.Second * 10)

			res, err := common.CheckIfTxMined(ip, rpc_port, txHash.(string))
			if err != nil {
				fmt.Printf("Failed to get tx's status, %v", err)
				os.Exit(-1)
			}
			if res != nil {
				resMap := res.(map[string]interface{})
				contractAddr := resMap["contractAddress"].(string)
				fmt.Printf("Congratulations! tx has been mined, contract address: %s \r\n", contractAddr)
				os.Exit(0)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVar(&contract, "contract", "", "name of contract to be deployed")
}
