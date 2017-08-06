## mcat
mcat是一个基于Golang实现的以太坊智能合约开发脚手架，它可以帮助你快速开发、调试以及部署智能合约，同时mact提供一个通过合约交易参数计算调用字节码的功能，可以帮助开发者不受语言限制，无论是Java, Python或是其它支持网络编程的语言与合约进行交互。

### 安装
#### 安装Solidity
Ubuntu下可通过以下方式快速安装`Solidity`，更多安装方式请查看[Solidity官方文档](https://solidity.readthedocs.io/en/develop/installing-solidity.html)
Ubuntu
```
sudo add-apt-repository ppa:ethereum/ethereum
sudo apt-get update
sudo apt-get install solc
```

#### 安装mcat
如果已经安装配置好`golang`，可以使用`golang`提供的工具快速安装
```
luren5@ubuntu:~$ go get github.com/luren5/mcat
luren5@ubuntu:~$ go install github.com/luren5/mcat
```
如果对`golang`不熟悉，可以直接下载编译好的安装包，放到`$PATHA`目录下

#### 检查安装是否成功
命令行下执行`mcat`，检查是否安装成功
```
luren5@ubuntu:~$ mcat
mcat is a development and testing framework for Ethereum implemented through golang.

Usage:
  mcat [command]

Available Commands:
  IDE         Solidity local online IDE.
  call        Call contract function.
  compile     compile contract
  deploy      deploy contract
  gasPrice    Show the current gas price.
  help        Help about any command
  init        init a new mc project
  loadConfig  Load config from the config file mcat.yaml.
  serve       Call contract function.

Flags:
  -h, --help   help for mcat

Use "mcat [command] --help" for more information about a command.
```

### mcat使用方法
注意：所有的mcat操作命令都必须在项目根目录下执行
#### 初始化项目
使用 `mcat init`初始化一个新的项目，初始化项目前请确定已经安装`git`，并且对当前目录有写权限，否则可能会造成初始化失败。

参数列表：
- `--project` 项目名，也是项目的目录名

使用示例：
```
luren5@ubuntu:~$ mcat init --project=mcat-demo
Congratulations! You have succeed in initing a new mc project.
```
`init`命令会clone[mcat模板项目](https://github.com/luren5/mcat-demo),  以下是各目录结构信息
```
- IDE                 # Solidity编辑器相关静态文件
- compiled            # 存放合约编译后的abi和bin文件
- contracts           # 合约源文件目录 
- data                # 配置文件加载后及项目相关数据存放目录
- mcat.yaml           # 配置文件
```

#### 加载项目配置
`mcat.yaml`为项目配置文件，可以根据需求配置多种模式下的不同配置，默认有`Development` 和 `Production` 两种
```
    project_name: "PROJECT_NAME"        # init操作会自动替换为`--project`参数值
    ip: "localhost"                     # 节点ip
    rpc_port: "8080"                    # 节点的rpc_port
    account: "0x34851ee7379fd43be25df08ab84b7402269fefc8"   # 默认发交易的账户，必须在你的节点中存在
    password: "123456"    # 密码账户的密码
    ide_port: "50728"     # IDE运行的端口
    server_port: "50729"  # serve服务运行的端口
```
使用`mcat loadConfig`来加载配置文件
参数列表：
- `model` 指定模式,  缺省情况下为`Development`

使用示例：
```
luren5@ubuntu:~/mcat-demo$ mcat loadConfig --model=Production
Succeed in loading mcat config.
```

#### 编译合约
使用`mcat compile`命令编译合约源码，合约源码放在`contracts`下,  比如示例项目中的`contracts/Ballot.sol`

参数列表：
- `sol` 指定需要编译的合约文件名
- `ext` 当指定的合约源文件中有多个合约时，`ext`参数可以排除不需要编译的合约, 多个用`,`隔开，编译后的合约字节码及`abi`存在在`compiled`目录下

使用示例：
```
luren5@ubuntu:~/mcat-demo$ mcat compile --sol=Ballot.sol --exc=Test
Waiting for compiling contracts…
Succeed in compiling contract Ballot
```

#### 部署合约
使用`mcat deploy命令部署合约`，部署合约前需要先启动节点，并且已完成配置加载（loadConfig）操作
参数列表：
- `--contract`  需要部署的合约

使用示例：
```
luren5@ubuntu:~/mcat-demo$ mcat deploy --contract=Ballot
Succeed in deploying contract Ballot, tx hash: 0xece1c87541d5cb16d81fa7fd055fa1bf259fac1e70221c2d5841145fed9e6474. Waiting for being mined…

Congratulations! tx has been mined, contract address: 0x0cf7ecedc79d011617da8144886c33caa799daa7
```
合约部署后，会返回交易的哈希值，`mcat`会一直等待交易被打包，每5s检查一次，直到交易被打包，返回合约账户地址

#### 调用合约方法
调用合约方法前需要先编译部署，并得到合约账户的地址
参数列表：
-  `--contract` 被调用的合约名
-  `--addr` 合约账户地址
-  `--function` 被调用的合约方法名
-  `--params` 参数列表，多个参数之间用`&`隔开，如果是数组类型参数，成员之间用`，`隔开
使用示例：
```
luren5@ubuntu:~/mcat-demo$ mcat call --contract Test2 --addr 0x23491c9c1b74bb15d988c495f945ec6e1b2720c2  --function cc --params 69

Succeed in calling cc, tx hash: 0xa15e4948a796136d9b461ce86b78f42b6a6ee13edd0aa00196a82e314641b63b
```
#### 获取gasPirce
获取所连接节点的当前gas price
```
luren5@ubuntu:~/mcat-demo$ mcat gasPrice
The current gas price is 0x4a817c800
```

#### 打开本地IDE
mcat提供一个简洁的本地IDE
使用示例
```
luren5@ubuntu:~/mcat-demo$ mcat IDE
IDE has been started, http://localhost:50728
```
启动后打开`http://localhost:50728` 即可使用IDE

#### 打开mcat server
mcat server是一个将调用合约方法相关参数 转换为交易的字节码（data）的服务器，这样开发者无论是使用`PHP`、`Python`、`Java`或者其它网络编程语言，只要遵循mcat的rpc规范，均可获得合约交易的data字段值，进而与合约进行交互

使用之前需要对合约进行编译，确实`compiled`目录下有与合约对应的`abi`
使用示例，另外需要启动`geth`节点以计算该合约调用需要消耗的`gas`
```
luren5@ubuntu:~/mcat-demo$ mcat serve
server has been started, listening 50727
# 新启一个终端窗口
luren5@ubuntu:~/mcat$ curl -X POST --data '{"jsonrpc":"2.0","method":"TxData.Detail","params":[{"Contract":"Ballot", "Function": "vote", "Params":"2"}],"id":67}' http://localhost:50727
{"id":67,"result":"{\"Bin\":\"0x0121b93f0000000000000000000000000000000000000000000000000000000000000002\",\"Gas\":\"0x53d9\"}","error":null}
``` 
