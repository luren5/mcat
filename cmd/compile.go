package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/luren5/mcat/utils"

	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/spf13/cobra"
)

var (
	sol  string
	solc string
	exc  string
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
			abi, _ := json.Marshal(contract.Info.AbiDefinition)
			//abi := string(abiBytes)
			bin := []byte(contract.Code)

			if err := ioutil.WriteFile(compiledDir+"/"+contractName+".abi", abi, 0660); err != nil {
				fmt.Printf("Failed to write compiling bin,  %v \r\n", err)
				os.Exit(-1)
			}

			if err := ioutil.WriteFile(compiledDir+"/"+contractName+".bin", bin, 0660); err != nil {
				fmt.Printf("Failed to write compiling bin,  %v \r\n", err)
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
	compileCmd.Flags().StringVar(&solc, "solc", "", "The path to solidity compiler.")

}
