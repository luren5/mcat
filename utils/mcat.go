package utils

import (
	"errors"
	"fmt"
	"os"

	"github.com/luren5/mcat/db"
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

func Config(key string) (string, error) {
	mDB, err := db.NewDB(db.DefaultPath)
	if err != nil {
		return "", err
	}
	v, err := mDB.Get([]byte(key))
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func GetDefaultAccount() (string, error) {
	if account, err := Config("account"); err != nil {
		return "", err
	} else {
		// if password was set, unlock the default account
		if password, err := Config("password"); err == nil {
			if ip, rpc_port, err := GetRpcInfo(); err == nil {
				params := fmt.Sprintf(`"%s", "%s"`, account, password)
				_, err = JrpcPost(ip, rpc_port, "personal_unlockAccount", params)
				if err != nil {
					return "", errors.New(fmt.Sprintf("Failed to unlock the default account, %v ", err))
				}

			}
		}

		return account, nil
	}
}

func GetRpcInfo() (string, string, error) {
	var ip, rpc_port string
	var err error
	if ip, err = Config("ip"); err != nil {
		errMes := fmt.Sprintf("Failed to read ip config, %v", err)
		return "", "", errors.New(errMes)
	}
	if rpc_port, err = Config("rpc_port"); err != nil {
		errMes := fmt.Sprintf("Failed to read rpc port config, %v", err)
		return "", "", errors.New(errMes)
	}
	return ip, rpc_port, nil
}
