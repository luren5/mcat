package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var project string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init a new mc project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(project) == 0 {
			fmt.Printf("Invalid project name \r\n")
			os.Exit(0)
		}
		fmt.Println("Starting the initialization, wait a moment…")

		demoUrl := "https://github.com/luren5/mcat-demo.git"
		if _, err := exec.Command("git", "clone", demoUrl, project).Output(); err != nil {
			fmt.Printf("Failed to clone demo project, %v \r\n", err)
			os.Exit(0)
		}

		configPath := project + "/mcat.yaml"
		buf, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic(err)
		}
		newContent := strings.Replace(string(buf), "PROJECT_NAME", project, -1)
		ioutil.WriteFile(configPath, []byte(newContent), 0)

		fmt.Println("Congratulations! You have succeed in initing a new mc project.")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&project, "project", "", "Project name.")
}
