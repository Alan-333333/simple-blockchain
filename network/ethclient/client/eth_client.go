package ethclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

// EthClient wraps ethereum client for convenient use.
type EthClient struct {
	client *ethclient.Client
	url    string // ethereum node URL
}

// NewEthClient creates a new EthClient instance.
func NewEthClient(url string) (*EthClient, error) {

	// Initialize geth client
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

// GetLatestBlockNumber retrieves the block number of the latest block.
func (ec *EthClient) GetLatestBlockNumber() (uint64, error) {

	header, err := ec.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}

	return header.Number.Uint64(), nil
}

// GetBalance retrieves the balance of the given account address.
func (ec *EthClient) GetBalance(account common.Address) (*big.Int, error) {

	balance, err := ec.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}

	return balance, nil

}

// SendTransaction sends the given transaction to blockchain.
func (ec *EthClient) SendTransaction(tx *types.Transaction) error {

	err := ec.client.SendTransaction(context.Background(), tx)
	if err != nil {
		return err
	}

	return nil

}

// GetTransactionByHash returns the transaction for the given hash.
func (ec *EthClient) GetTransactionByHash(hash common.Hash) (*types.Transaction, bool, error) {

	tx, pending, err := ec.client.TransactionByHash(context.Background(), hash)
	if err != nil {
		return nil, false, err
	}

	return tx, pending, nil
}

// GetTransactionReceipt returns the receipt for the given transaction hash.
func (ec *EthClient) GetTransactionReceipt(hash common.Hash) (*types.Receipt, error) {

	receipt, err := ec.client.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

// FilterLogs filters logs based on the given query.
func (ec *EthClient) FilterLogs(query ethereum.FilterQuery) ([]types.Log, error) {

	logs, err := ec.client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// rpcCall sends an JSON RPC request and parses the response.
func (ec *EthClient) rpcCall(
	method string,
	args []interface{},
	result interface{},
) error {

	reqID := uuid.New()
	// Construct request
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  args,
		"id":      reqID,
	}

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	// Send request
	resp, err := http.Post(ec.url, "application/json", bytes.NewReader(reqData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Parse response
	var respBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return err
	}

	error := respBody["error"]
	if error != nil {
		return fmt.Errorf("%v", error)
	}

	// Set result
	result = respBody["result"]

	return nil
}
