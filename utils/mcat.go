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

func GetDefaultAccount() (string, error) {
	var account string
	if r, err := Config("account"); err != nil {
		return "", errors.New(fmt.Sprintf("Invalid default account, %v", err))
	} else {
		account = r.(string)
	}

	if r, err := Config("password"); err != nil {
		fmt.Printf("You haven't set password of default account, please make sure the default account has been unlocked.")
	} else {
		password := r.(string)
		// unlock
		ip, rpc_port, err := GetRpcInfo()
		if err != nil {
			return "", err
		}

		params := fmt.Sprintf(`"%s", "%s"`, account, password)
		_, err = JrpcPost(ip, rpc_port, "personal_unlockAccount", params)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Failed to unlock the default account, %v ", err))
		}
	}
	return account, nil
}

func GetRpcInfo() (string, string, error) {
	var ip, rpc_port string
	if r, err := Config("ip"); err != nil {
		errMes := fmt.Sprintf("Failed to read ip config, %v", err)
		return "", "", errors.New(errMes)
	} else {
		ip = r.(string)
	}
	if r, err := Config("rpc_port"); err != nil {
		errMes := fmt.Sprintf("Failed to read rpc port config, %v", err)
		return "", "", errors.New(errMes)
	} else {
		rpc_port = r.(string)
	}
	return ip, rpc_port, nil
}
