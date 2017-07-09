package ethabi

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

type EthABI struct {
	Contract string
	ABI      []byte
}

type ContractParam struct {
	Type       string
	Value      interface{}
	MemberType uint
}

const (
	MemberTypeS = iota
	MemberTypeI
)

func NewEthABI(contract string, abi []byte) *EthABI {
	return &EthABI{
		Contract: contract,
		ABI:      abi,
	}
}

func CalBytes(selector string, cp []ContractParam) (string, error) {
	var paramData []string
	var dynaData []string
	var paramNum = len(cp)
	var offset = paramNum * 32

	for _, v := range cp {
		// bytes<M>
		if m, _ := regexp.MatchString(`^bytes\d+$`, v.Type); m {
			paramData = append(paramData, fmt.Sprintf("%x", v.Value.(string))+strings.Repeat("0", 64-2*len(v.Value.(string))))
			continue
		}

		/********** the following are dynamic ************/
		// string or bytes
		if v.Type == "string" || v.Type == "bytes" {
			paramData = append(paramData, fmt.Sprintf("%064x", offset))
			offset = offset + 2*32 // len + content

			val := v.Value.(string)
			dynaData = append(dynaData, fmt.Sprintf("%064x", len(val)))
			dynaData = append(dynaData, fmt.Sprintf("%x", val)+strings.Repeat("0", 64-2*len(val)))
			continue
		}

		// T[] or T[k]
		if m, _ := regexp.MatchString(`^.*\[\d?\]$`, v.Type); m {
			paramData = append(paramData, fmt.Sprintf("%064x", offset))

			val := v.Value.([]string)
			offset = offset + (1+len(val))*32 // num of + every data
			dynaData = append(dynaData, fmt.Sprintf("%064x", len(val)))

			if v.MemberType == MemberTypeS { // string or bytes
				for _, p := range val {
					dynaData = append(dynaData, fmt.Sprintf("%x", p)+strings.Repeat("0", 64-2*len(p)))
				}
			} else { // int
				for _, p := range val {
					pi, err := strconv.Atoi(p)
					if err != nil {
						errMes := fmt.Sprintf("Invalid %s value: %s", v.Type, p)
						return "", errors.New(errMes)
					}
					dynaData = append(dynaData, fmt.Sprintf("%064x", pi))
				}
			}
			continue
		}

		/*********************** default ********************/
		paramData = append(paramData, fmt.Sprintf("%064x", v.Value))
	}

	/*
		fmt.Println(paramData)
		fmt.Println(dynaData)
	*/
	return "0x" + selector + strings.Join(paramData, "") + strings.Join(dynaData, ""), nil
}

func (e *EthABI) FuncDef(function string) (string, error) {
	var aa []map[string]interface{}
	json.Unmarshal(e.ABI, &aa)

	for _, v := range aa {
		if v["name"] == function {
			var paramTypes []string
			inputs := v["inputs"].([]interface{})

			for _, p := range inputs {
				input := p.(map[string]interface{})
				t := input["type"].(string)

				// T[k]
				if m, _ := regexp.MatchString(`^.*\[\d+\]$`, t); m {
					ptn := `^(\w+?)(\d+)\[(\d+)\]$`
					reg := regexp.MustCompile(ptn)
					r := reg.FindSubmatch([]byte(t))

					st := string(r[1])
					l := string(r[2])
					n := string(r[3])
					paramTypes = append(paramTypes, fmt.Sprintf("%s%sx%s[%s]", st, l, l, n))
					continue
				}
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
	if len(paramSlice) == 0 {
		return nil, nil
	}

	var aa []map[string]interface{}
	json.Unmarshal(e.ABI, &aa)

	var inputs []interface{}
	for _, v := range aa {
		if v["name"].(string) == function {
			inputs = v["inputs"].([]interface{})
			break
		}
	}
	if len(inputs) != len(paramSlice) {
		errMes := fmt.Sprintf("Params num does't match, expecting %d and actually %d", len(inputs), len(paramSlice))
		return nil, errors.New(errMes)
	}

	var cp []ContractParam
	for i, p := range inputs {
		si := p.(map[string]interface{})
		siType := si["type"].(string)
		siValue := paramSlice[i]

		// T[k]
		if m, _ := regexp.MatchString(`^.*\[\d+\]$`, siType); m {
			var memberType uint
			if strings.Contains(siType, "int") {
				memberType = MemberTypeI
			} else {
				memberType = MemberTypeS
			}

			valSlice := strings.Split(siValue, ",")
			c := ContractParam{Type: siType, Value: valSlice, MemberType: memberType}
			cp = append(cp, c)
			continue
		}

		c := ContractParam{Type: siType, Value: siValue}
		cp = append(cp, c)
	}

	fmt.Println(cp) //, inputs)
	return cp, nil
}

/*
func typeRepresent(t string) string {
	switch t {
	case "uint":
		return "uint256"
	case "int":
		return "int256"
	case "fixed":
		return "fixed128"
	case "ufixed":
		return "unfixed128"
	}
	return t
}
*/