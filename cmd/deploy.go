// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	yaml "github.com/ghodss/yaml"

	"github.com/luren5/mcat/common"
	"github.com/luren5/mcat/utils"
	"github.com/spf13/cobra"
)

var (
	contract string
	ip       string
	rpc_port string
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy contract",
	Long:  `create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// read bin
		compiledFile := utils.CompiledDir() + contract
		data, err := ioutil.ReadFile(compiledFile)
		if err != nil {
			fmt.Printf("Failed to read contract file, %v", err)
			os.Exit(-1)
			return
		}
		var cc common.Compiled
		yaml.Unmarshal(data, &cc)

		// read config
		data, err = ioutil.ReadFile(utils.ConfigPath())
		if err != nil {
			fmt.Printf("Failed to read config file, %v\r\n", err)
			os.Exit(-1)
		}
		var dm map[string]interface{}
		yaml.Unmarshal(data, &dm)

		var config map[string]interface{}
		if v, ok := dm[dm["model"].(string)]; ok {
			config = v.(map[string]interface{})
		} else {
			fmt.Sprintf("Failed to get config, %v", err)
			os.Exit(-1)
		}
		ip = config["ip"].(string)
		rpc_port = config["rpc_port"].(string)

		if config["ip"] == nil || len(ip) == 0 {
			fmt.Sprintf("Invalid ip config.")
			os.Exit(-1)
		}
		if config["rpc_port"] == nil || len(rpc_port) == 0 {
			fmt.Sprintf("Invalid rpc port config.")
			os.Exit(-1)
		}

		if config["account"] == nil || len(config["account"].(string)) != 42 || !strings.HasPrefix(config["account"].(string), "0x") {
			fmt.Println(len(config["account"].(string)))
			os.Exit(-1)
		}

		if config["password"] == nil || len(config["password"].(string)) == 0 {
			fmt.Printf("You haven't set password of default account, please make sure the default account has been unlocked.")
		} else {
			// unlock
			params := fmt.Sprintf(`"%s", "%s"`, config["account"], config["password"])
			_, err := utils.JrpcPost(ip, rpc_port, "personal_unlockAccount", params)
			if err != nil {
				fmt.Printf("Failed to unlock the default account, %v ", err)
				os.Exit(-1)
			}
		}

		// cal gas
		gas := "0x76ffce" // here needs cal
		gasPrice := "0x9184e72a000"
		// do deploy
		tx := new(common.Transaction)
		tx.From = config["account"].(string)
		tx.Gas = gas
		tx.GasPrice = gasPrice
		tx.Data = cc.Bin
		tx.Type = common.TxTypeContract

		txHash, err := utils.SendTransaction(ip, rpc_port, tx)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		fmt.Printf("Succeed in deploying contract %s, tx hash: %s. Waiting for being mined…\r\n \r\n", contract, txHash)

		// check status
		for {
			time.Sleep(time.Second * 10)

			res, err := utils.CheckIfTxMined(ip, rpc_port, txHash.(string))
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
