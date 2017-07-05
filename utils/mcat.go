package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var compiledDir string = "compiled"
var contractsDir string = "contracts"

func ContractsDir() string {
	pwd, _ := os.Getwd()
	return pwd + "/" + contractsDir + "/"
}

func CompiledDir() string {
	pwd, _ := os.Getwd()
	return pwd + "/" + compiledDir + "/"
}

func ConfigPath() string {
	pwd, _ := os.Getwd()
	return pwd + "/mcat.yaml"
}

func Config(key string) (interface{}, error) {
	keys := strings.Split(key, ".")
	configPath := ConfigPath()
	if _, err := os.Stat(configPath); err != nil {
		return nil, err
	}
	// read config
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var dm map[string]interface{}
	yaml.Unmarshal(data, &dm)

	model := dm["model"].(string)
	config, ok := dm[model] // develop or product
	if !ok {
		return nil, errors.New(fmt.Sprintf("config %s not found", model))
	}

	var res map[interface{}]interface{} = config.(map[interface{}]interface{})
	kl := len(keys)
	for i := 0; i < kl; i++ {
		key := interface{}(keys[i])
		c, ok := res[key]
		if !ok {
			return nil, errors.New(fmt.Sprintf("config %s not found", key.(string)))
		}
		if i == kl-1 { // end
			return c, nil
		}
		res = c.(map[interface{}]interface{})
	}

	return res, nil
}
