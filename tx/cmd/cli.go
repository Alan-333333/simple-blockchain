package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	blockchain "github.com/Alan-333333/simple-blockchain/block/chain"
	"github.com/Alan-333333/simple-blockchain/tx"
	"github.com/Alan-333333/simple-blockchain/utils"
)

func main() {
	// 1. 生成密钥对
	privKey, pubKey := utils.GenerateKeyPair()

	// 2. 创建交易
	txs := tx.NewTransaction(pubKey, "receiver", 10)

	// 3. 签名交易
	txs.Sign(privKey)

	// 4. 创建交易池
	pool := tx.NewTxPool()
	pool.AddTx(txs)

	// 5. 生成创世区块
	genesisBlock := &blockchain.Block{
		Version:   blockchain.CURRENT_BLOCK_VERSION,
		Timestamp: uint64(time.Now().Unix()),
	}
	genesisBlock.Transactions = pool.Txs

	// 6. 持久化
	saveBlock(genesisBlock)

	fmt.Println("创世区块生成完成!")
}

func saveBlock(block *blockchain.Block) {
	// 省略区块保存实现
	// 2. 序列化
	blockData, _ := json.Marshal(block)
	// 3. 存储到文件
	err := os.WriteFile("genesis.blk", blockData, 0644)
	if err != nil {
		// 错误处理
	}
	// 4. 读取并反序列化
	data, _ := os.ReadFile("genesis.blk")
	var savedBlock blockchain.Block
	json.Unmarshal(data, &savedBlock)
	// 5. 验证
	if reflect.DeepEqual(block, savedBlock) {
		fmt.Println("创世区块保存成功!")
	}
}
