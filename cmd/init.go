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
	"os"

	"github.com/spf13/cobra"
)

var projectName string = "testcba"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init a new mc project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := os.Mkdir(projectName, 0777); err != nil {
			fmt.Printf("Failed to make dir, %v \r\n", err)
			os.Exit(0)
		}
		if err := os.Chdir("testcba"); err != nil {
			fmt.Printf("Failed change dir, %v \r\n", err)
			os.Exit(0)
		}

		fmt.Println(os.Getwd())

		// download demo

		fmt.Println("Congratulations! You have succeed in initing a new mc project.")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

}
