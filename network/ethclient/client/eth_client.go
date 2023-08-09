package ethclient

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthClient 封装以太坊客户端
type EthClient struct {
	client *ethclient.Client
	url    string // 节点URL
}

// NewEthClient 创建客户端
func NewEthClient(url string) (*EthClient, error) {
	// 初始化geth客户端
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	ec := &EthClient{
		client: client,
		url:    url,
	}
	return ec, nil
}

// GetLatestBlock 获取最新区块
func (ec *EthClient) GetLatestBlock() (string, error) {
	header, err := ec.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return "", err
	}

	return header.Number.String(), nil
}

func (c *EthClient) GetBalance(account common.Address) (*big.Int, error) {
	return c.client.BalanceAt(context.Background(), account, nil)
}

func (c *EthClient) SendTransaction(tx *types.Transaction) error {
	return c.client.SendTransaction(context.Background(), tx)
}

func (c *EthClient) GetTransaction(txHash common.Hash) (*types.Transaction, bool, error) {
	return c.client.TransactionByHash(context.Background(), txHash)
}

func (c *EthClient) GetTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	return c.client.TransactionReceipt(context.Background(), txHash)
}

func (c *EthClient) GetLogs(query ethereum.FilterQuery) ([]types.Log, error) {
	return c.client.FilterLogs(context.Background(), query)
}

// EthClient中的rpcCall方法

func (ec *EthClient) rpcCall(
	ctx context.Context,
	result interface{},
	method string,
	args ...interface{},
) error {

	// 构造请求
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  args,
		"id":      123,
	}
	data, _ := json.Marshal(req)

	// 发送请求
	r, err := http.Post(ec.url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// 解析响应
	var resp map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return err
	}
	error := resp["error"]
	if error != nil {
		return err
	}

	// 设置结果
	resultVal := resp["result"]

	switch r := result.(type) {
	case *big.Int:
		*r = *resultVal.(*big.Int)
	case *types.Block:
		*r = *resultVal.(*types.Block)
	}
	return nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}
