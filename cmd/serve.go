package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"net/http"

	"github.com/luren5/mcat/common"
	"github.com/luren5/mcat/ethabi"
	"github.com/luren5/mcat/utils"
	"github.com/spf13/cobra"
)

type HttpConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

type TxData struct {
	Bin string
	Gas string
}

type Args struct {
	Contract, Function, Params string
}

func (t *TxData) Detail(args *Args, reply *string) error {
	bin, err := getBin(args)
	if err != nil {
		return err
	}

	ip, rpc_port, err := utils.GetRpcInfo()
	if err != nil {
		return err
	}

	tx := new(common.Transaction)
	tx.From = common.ZeroAddr
	tx.To = common.ZeroAddr
	tx.Data = bin
	tx.Type = common.TxTypeCommon

	if gas, err := common.EstimateGas(ip, rpc_port, tx); err != nil {
		return err
	} else {
		td := new(TxData)
		td.Bin = bin
		td.Gas = gas

		r, _ := json.Marshal(td)
		*reply = string(r)
		return nil
	}
}

func getBin(args *Args) (string, error) {
	if len(args.Contract) == 0 {
		return "", errors.New("Invalid contract name")
	}
	if len(args.Function) == 0 {
		return "", errors.New("Invalid function name")
	}
	// ethabi
	e, err := ethabi.NewEthABI(args.Contract)
	if err != nil {
		return "", err
	}

	// selector
	funcDef, err := e.FuncDef(args.Function)
	if err != nil {
		return "", err
	}
	selector := ethabi.CalSelector(funcDef)

	// params
	cp, err := e.CombineParams(args.Function, args.Params)
	if err != nil {
		return "", err
	}

	callBytes, err := ethabi.CalBytes(selector, cp)
	if err != nil {
		return "", err
	}

	return callBytes, nil
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Call contract function.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// rpc server
		var server = rpc.NewServer()
		server.Register(new(TxData))
		var server_port string
		if s, err := utils.Config("server_port"); err != nil {
			server_port = "50729"
		} else {
			server_port = s.(string)
		}

		listener, err := net.Listen("tcp", ":"+server_port)
		if err != nil {
			fmt.Println("listen error:", err)
		}
		defer listener.Close()
		fmt.Println("server has been started, listening " + server_port)

		http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				serverCodec := jsonrpc.NewServerCodec(&HttpConn{in: r.Body, out: w})
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(200)
				err := server.ServeRequest(serverCodec)
				if err != nil {
					log.Printf("Error while serving JSON request: %v", err)
					http.Error(w, "Error while serving JSON request, details have been logged.", 500)
					return
				}

			}
		}))
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
