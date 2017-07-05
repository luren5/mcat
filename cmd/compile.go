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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "github.com/ghodss/yaml"
	"github.com/luren5/mcat/common"
	"github.com/luren5/mcat/utils"

	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/spf13/cobra"
)

var (
	sol string
	exc string
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "compile contract",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(sol) == 0 {
			fmt.Println("No input file specified.")
			os.Exit(-1)
		}
		exclude := make(map[string]bool)
		for _, kind := range strings.Split(exc, ",") {
			exclude[kind] = true
		}
		contracts, err := compiler.CompileSolidity("solc", utils.ContractsDir()+sol)

		fmt.Println("Waiting for compiling contracts…")

		if err != nil {
			fmt.Printf("Failed to compile contract, %v \r\n", err)
			os.Exit(-1)
		}

		compiledDir := utils.CompiledDir()
		if _, err := os.Stat(compiledDir); err != nil {
			os.MkdirAll(compiledDir, 0777)
		}

		for name, contract := range contracts {
			nameParts := strings.Split(name, ":")
			contractName := nameParts[len(nameParts)-1]
			if exclude[contractName] {
				continue
			}
			abiBytes, _ := json.Marshal(contract.Info.AbiDefinition)
			abi := string(abiBytes)
			bin := contract.Code

			compiled := new(common.Compiled)
			compiled.Name = name
			compiled.Abi = abi
			compiled.Bin = bin
			compiledContent, _ := yaml.Marshal(compiled)

			if err := ioutil.WriteFile(compiledDir+"/"+contractName, compiledContent, 0660); err != nil {
				fmt.Printf("Failed to write compiling,  %v \r\n", err)
				os.Exit(-1)
			}
			fmt.Printf("Succeed in compiling contract %s \r\n", contractName)
		}

	},
}

func init() {
	RootCmd.AddCommand(compileCmd)

	compileCmd.Flags().StringVar(&sol, "sol", "", "Path to contract source file to compile.")
	compileCmd.Flags().StringVar(&exc, "exc", "", "Comma separated types to exclude from compiling.")
}
