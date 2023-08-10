package main

import (
	"fmt"
	"time"

	blockchain "github.com/Alan-333333/simple-blockchain/block/chain"
	transaction "github.com/Alan-333333/simple-blockchain/transaction"
	"github.com/Alan-333333/simple-blockchain/utils"
)

func main() {
	// 1. 生成密钥对
	privKey, pubKey := utils.GenerateKeyPair()

	// 2. 创建交易
	address := utils.PubKeyToAddress(pubKey)

	txs := transaction.NewTransaction(address, "receiver", 10)

	// 3. 签名交易
	txs.Sign(privKey)

	// 4. 创建交易池
	pool := transaction.NewTxPool()
	pool.AddTx(txs)

	// 5. 生成创世区块
	genesisBlock := &blockchain.Block{
		Version:   blockchain.CURRENT_BLOCK_VERSION,
		Timestamp: uint64(time.Now().Unix()),
	}
	genesisBlock.Transactions = pool.Txs

	// 6. 持久化
	// pool.Save()

	fmt.Println("创世区块生成完成!")
}
