## mcat
一个基于Golang实现的以太坊智能合约开发框架

### 安装
#### Ubuntu
如果已经安装配置好`golang`，可以使用`golang`提供的工具快速安装
```
luren5@ubuntu:~$ go get github.com/luren5/mcat
luren5@ubuntu:~$ go install github.com/luren5/mcat
```


#### 检查安装是否成功
安装完成后，在命令行下执行`mcat`，如果能正常
```
luren5@ubuntu:~$ mcat
mcat is a development and testing framework for Ethereum implemented through golang.

Usage:
  mcat [command]

Available Commands:
  compile     compile contract
  deploy      deploy contract
  help        Help about any command
  init        init a new mc project

Flags:
      --config string   config file (default is $HOME/.mcat.yaml)
  -h, --help            help for mcat
  -t, --toggle          Help message for toggle

Use "mcat [command] --help" for more information about a command.
```

### 使用方法
#### 初始化项目
使用 `mcat init`初始化一个新的项目，初始化项目前请确定已经安装`git`，并且对当前目录有写权限，否则可能会造成初始化失败。

参数列表：
- `--project` 项目名，也是项目的目录名

使用示例：
```
luren5@ubuntu:~$ mcat init --project=mcat-demo
Congratulations! You have succeed in initing a new mc project.
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
使用`mcat deploy命令部署合约`，部署合约前需要先启动节点，并修改配置文件`mcat.yaml`中的默认配置
```
model: "development"      # 当前开发模式

development:
    ip: "localhost"       # 节点的IP
    rpc_port: "8080"      # 节点的RPC服务端口
    account: "0x34851ee7379fd43be25df08ab84b7402269fefc8"  # 默认用来发交易的账户
    password: "123456"    # 默认账户的密码
```

参数列表：
- `--contract`  需要部署的合约

使用示例：
```
luren5@ubuntu:~/mcat-demo$ mcat deploy --contract=Ballot
Succeed in deploying contract Ballot, tx hash: 0xece1c87541d5cb16d81fa7fd055fa1bf259fac1e70221c2d5841145fed9e6474. Waiting for being mined…

Congratulations! tx has been mined, contract address: 0x0cf7ecedc79d011617da8144886c33caa799daa7
```
合约部署后，会返回交易的哈希值，`mcat`会一直等待交易被打包，每5s检查一次，直到交易被打包，返回合约账户地址
