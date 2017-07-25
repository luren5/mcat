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

func LastEditFile(args ...interface{}) (string, error) {
	mDB, err := db.NewDB(db.DefaultPath)
	if err != nil {
		return "", err
	}
	key := []byte("lastEditFile")
	if len(args) == 0 { // get
		if val, err := mDB.Get(key); err != nil {
			return "", err
		} else {
			return string(val), nil
		}
	} else { // set
		f := args[0]
		fileName := f.(string)
		err := mDB.Put(key, []byte(fileName))
		return "", err
	}
}
func Config(key string) (interface{}, error) {
	mDB, err := db.NewDB(db.DefaultPath)
	if err != nil {
		return nil, err
	}
	val, err := mDB.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	return string(val), nil
}

func ProjectName() string {
	r, err := Config("project_name")
	if err != nil {
		return ""
	}
	return r.(string)
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
