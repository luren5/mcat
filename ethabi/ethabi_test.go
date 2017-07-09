package ethabi

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestCallBaz(t *testing.T) {
	contract := "Test"
	function := "baz"
	params := "69&true"
	res := "0xcdcd77c000000000000000000000000000000000000000000000000000000000000000450000000000000000000000000000000000000000000000000000000000000001"

	abiFile := "../compiled/Test.abi"
	abiBytes, err := ioutil.ReadFile(abiFile)
	if err != nil {
		panic(err)
	}

	e := EthABI{contract, abiBytes}
	def, err := e.FuncDef(function)
	if err != nil {
		panic(err)
	}
	selector := CalSelector(def)
	cp, err := e.CombineParams(function, params)
	if err != nil {
		panic(err)
	}
	r, err := CalBytes(selector, cp)
	if err != nil {
		panic(err)
	}
	if r != res {
		fmt.Println("not equal")
	}
}

func TestCallF(t *testing.T) {
	result := "0x8be6524600000000000000000000000000000000000000000000000000000000000001230000000000000000000000000000000000000000000000000000000000000080313233343536373839300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000004560000000000000000000000000000000000000000000000000000000000000789000000000000000000000000000000000000000000000000000000000000000d48656c6c6f2c20776f726c642100000000000000000000000000000000000000"

	funcDef := "f(uint256,uint32[],bytes10,bytes)"
	fs := CalSelector(funcDef)
	p0 := ContractParam{Type: "uint", Value: "291"}
	p1 := ContractParam{"uint32[]", "1110,1929", MemberTypeI}
	p2 := ContractParam{Type: "bytes10", Value: "1234567890"}
	p3 := ContractParam{Type: "bytes", Value: "Hello, world!"}

	var cp = [...]ContractParam{p0, p1, p2, p3}
	cb, err := CalBytes(fs, cp[:])
	if err != nil {
		panic(err)
	}
	if cb != result {
		fmt.Println("not equal")
	}
}
