# EthClient

EthClient是一个简单、方便的以太坊Go语言客户端。

## 特性

- 提供从以太坊节点获取数据的方法,如区块、交易、账户余额等
- 封装了JSON-RPC请求和响应处理
- 简单的CLI交互接口,使得命令行使用方便

## 用法

### 作为库

```go
client, err := ethclient.NewEthClient(endpoint)
balance, err := client.GetBalance(account)
```

## 运行

需要安装Go 1.19+

```bash
$ go run cmd/cli.go
```

## 开发

```bash
# 下载依赖
$ go mod download

# 格式化代码
$ go fmt ./...

# 测试
$ go test ./...
```

欢迎贡献!

## 许可

Apache 2.0许可证
