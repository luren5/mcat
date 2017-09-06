package ethabi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/luren5/mcat/utils"
)

type EthABI struct {
	Contract string
	ABI      []byte
}

type ContractParam struct {
	Type       string
	Value      string
	MemberType uint
}

const (
	MemberTypeS = iota
	MemberTypeI
	MemberTypeF
)

func NewEthABI(contract string) (*EthABI, error) {
	abiFile := utils.CompiledDir() + contract + ".abi"
	abiBytes, err := ioutil.ReadFile(abiFile)
	if err != nil {
		return nil, err
	}
	return &EthABI{
		Contract: contract,
		ABI:      abiBytes,
	}, nil
}

func CalBytes(selector string, cp []ContractParam) (string, error) {
	var paramData []string
	var dynaData []string
	var paramNum = len(cp)
	var offset = paramNum * 32

	for _, v := range cp {
		// bytes<M>
		if m, _ := regexp.MatchString(`^bytes\d+$`, v.Type); m {
			paramData = append(paramData, fmt.Sprintf("%x", v.Value)+strings.Repeat("0", 64-2*len(v.Value)))
			continue
		}
		// int or uint
		if m, _ := regexp.MatchString(`int\d*$`, v.Type); m {
			bi := big.NewInt(0)
			if _, ok := bi.SetString(v.Value, 10); !ok {
				errMes := fmt.Sprintf("Failed to convert %s to big int", v.Value)
				return "", errors.New(errMes)
			}

			/*
				val, err := strconv.Atoi(v.Value)
				if err != nil {
					return "", err
				}
			*/
			fmt.Println("bi: ", bi)
			paramData = append(paramData, fmt.Sprintf("%064x", bi))
		}
		// fixed or ufixed
		if m, _ := regexp.MatchString(`fixed\d*$`, v.Type); m {
			//fmt.Println(v.Type)
		}

		/********** the following are dynamic ************/
		// string or bytes
		if v.Type == "string" || v.Type == "bytes" {
			paramData = append(paramData, fmt.Sprintf("%064x", offset))
			offset = offset + 2*32 // len + content

			val := v.Value
			dynaData = append(dynaData, fmt.Sprintf("%064x", len(val)))
			dynaData = append(dynaData, fmt.Sprintf("%x", val)+strings.Repeat("0", 64-2*len(val)))
			continue
		}

		// T[] or T[k]
		if m, _ := regexp.MatchString(`^.*\[\d*\]$`, v.Type); m {
			paramData = append(paramData, fmt.Sprintf("%064x", offset))

			val := strings.Split(v.Value, ",")
			offset = offset + (1+len(val))*32 // num of + every data
			dynaData = append(dynaData, fmt.Sprintf("%064x", len(val)))

			switch v.MemberType {
			case MemberTypeS:
				for _, p := range val {
					dynaData = append(dynaData, fmt.Sprintf("%x", p)+strings.Repeat("0", 64-2*len(p)))
				}
			case MemberTypeI:
				for _, p := range val {
					/*
						pi, err := strconv.Atoi(p)
								if err != nil {
									errMes := fmt.Sprintf("Invalid %s value: %s", v.Type, p)
									return "", errors.New(errMes)
								}
					*/

					bi := big.NewInt(0)
					if _, ok := bi.SetString(p, 10); !ok {
						errMes := fmt.Sprintf("Failed to convert %s to big int", v.Value)
						return "", errors.New(errMes)
					}
					dynaData = append(dynaData, fmt.Sprintf("%064x", bi))
				}
			case MemberTypeF: // fixed and ufixed hasn't been implemented yet

			}
			continue
		}
	}
	return "0x" + selector + strings.Join(paramData, "") + strings.Join(dynaData, ""), nil
}

func fixedBytes(pf float64) string {
	return ""
}

func (e *EthABI) FuncDef(function string) (string, error) {
	var aa []map[string]interface{}
	json.Unmarshal(e.ABI, &aa)

	for _, v := range aa {
		if v["name"].(string) == function {
			var paramTypes []string
			inputs := v["inputs"].([]interface{})

			for _, p := range inputs {
				input := p.(map[string]interface{})
				t := input["type"].(string)
				paramTypes = append(paramTypes, t)
			}
			def := fmt.Sprintf("%s(%s)", function, strings.Join(paramTypes, ","))
			return def, nil
		}
	}

	return "", errors.New("Can't find the specified function.")
}

func CalSelector(funcDef string) string {
	h := crypto.Keccak256Hash([]byte(funcDef))
	return fmt.Sprintf("%x", h[0:4])
}

func (e *EthABI) CombineParams(function, paramStr string) ([]ContractParam, error) {
	paramSlice := strings.Split(paramStr, "&")
	if len(paramStr) == 0 || len(paramSlice) == 0 {
		return nil, nil
	}

	var aa []map[string]interface{}
	json.Unmarshal(e.ABI, &aa)

	// get inputs
	var inputs []interface{}
	var target string
	if function == "constructor" {
		target = "type"
	} else {
		target = "name"
	}
	for _, v := range aa {
		if v[target].(string) == function {
			inputs = v["inputs"].([]interface{})
			break
		}
	}

	// check if num equal
	if len(inputs) != len(paramSlice) {
		errMes := fmt.Sprintf("Params num does't match, expecting %d and actually %d", len(inputs), len(paramSlice))
		return nil, errors.New(errMes)
	}

	var cp []ContractParam
	for i, p := range inputs {
		si := p.(map[string]interface{})
		siType := si["type"].(string)
		siValue := paramSlice[i]

		// bool
		if siType == "bool" {
			siType = "uint8"
			if siValue == "true" {
				siValue = "1"
			} else {
				siValue = "0"
			}
		}

		// T[k]
		if m, _ := regexp.MatchString(`^.*\[\d*\]$`, siType); m {
			var memberType uint
			if strings.Contains(siType, "int") {
				memberType = MemberTypeI
			}
			if strings.Contains(siType, "fixed") {
				memberType = MemberTypeF
			}
			if strings.Contains(siType, "bytes") || strings.Contains(siType, "string") {
				memberType = MemberTypeS
			}

			c := ContractParam{Type: siType, Value: siValue, MemberType: memberType}
			cp = append(cp, c)
			continue
		}

		c := ContractParam{Type: siType, Value: siValue}
		cp = append(cp, c)
	}

	return cp, nil
}
