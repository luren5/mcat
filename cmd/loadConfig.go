package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/luren5/mcat/db"
	"github.com/luren5/mcat/utils"
	"github.com/spf13/cobra"
)

var model string = "Development"

// loadConfigCmd represents the loadConfig command
var loadConfigCmd = &cobra.Command{
	Use:   "loadConfig",
	Short: "Load config from the config file mcat.yaml.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if model != "Production" {
			model = "Development"
		}
		configFile := utils.ConfigPath()
		if _, err := os.Stat(configFile); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		// read config
		data, err := ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Printf("Failed to read config, %v", err)
			os.Exit(-1)
		}

		var dm map[string]interface{}
		yaml.Unmarshal(data, &dm)
		config, ok := dm[model]
		if !ok {
			fmt.Printf("%s config not found", model)
			os.Exit(-1)
		}

		mDB, err := db.NewDB(db.DefaultPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}

		var configMap map[interface{}]interface{} = config.(map[interface{}]interface{})
		for k, v := range configMap {
			key := []byte(k.(string))
			val := []byte(v.(string))
			if len(key) > 0 && len(val) > 0 {
				err := mDB.Put(key, val)
				if err != nil {
					fmt.Printf("Failed to load config, key %s val %s, %v", k.(string), v.(string), err)
					os.Exit(-1)
				}
			}
		}
		fmt.Println("Succeed in loading mcat config.")
	},
}

func init() {
	RootCmd.AddCommand(loadConfigCmd)
	loadConfigCmd.Flags().StringVar(&model, "model", "", "Current development mode.")
}
