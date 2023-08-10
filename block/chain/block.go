package blockchain

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/Alan-333333/simple-blockchain/transaction"
)

const CURRENT_BLOCK_VERSION = 1

// Block结构体代表区块
type Block struct {
	// 版本号
	Version uint64

	// 前一个区块的Hash
	PrevHash []byte

	// Merkle树的根Hash
	MerkleRoot []byte

	// 当前区块创建的时间
	Timestamp uint64

	// 难度目标
	Difficulty uint64

	// 随机数,将与Nonce参与挖矿
	Nonce []byte

	// 当前区块的Hash
	Hash []byte

	// 该区块中的交易列表
	Transactions []*transaction.Transaction
}

func NewBlock() *Block {
	return &Block{
		Version:   CURRENT_BLOCK_VERSION,
		Timestamp: uint64(time.Now().Unix()),
	}
}
func (block *Block) Save() {
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
	var savedBlock Block
	json.Unmarshal(data, &savedBlock)
	// 5. 验证
	if reflect.DeepEqual(block, savedBlock) {
		fmt.Println("区块保存成功!")
	}
}
