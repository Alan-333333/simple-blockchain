package blockchain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/Alan-333333/simple-blockchain/transaction"
)

const CURRENT_BLOCK_VERSION = 1
const BASE_BLOCK_DIFFCULTY = 1

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

const genesisFile = "./dat/blockchain/genesis.blk"

func NewBlock(prevHash []byte, prevDiffculty uint64) *Block {
	return &Block{
		PrevHash:   prevHash,
		Version:    CURRENT_BLOCK_VERSION,
		Difficulty: prevDiffculty,
		Timestamp:  uint64(time.Now().Unix()),
	}
}

func (block *Block) Save() {
	// 省略区块保存实现
	// 2. 序列化
	blockData, _ := json.Marshal(block)
	// 3. 存储到文件
	err := os.WriteFile(genesisFile, blockData, 0644)
	if err != nil {
		// 错误处理
	}
	// 4. 读取并反序列化
	data, _ := os.ReadFile(genesisFile)
	var savedBlock Block
	json.Unmarshal(data, &savedBlock)
	// 5. 验证
	if reflect.DeepEqual(block, savedBlock) {
		fmt.Println("区块保存成功!")
	}
}

func CreateGenesisBlock() *Block {
	// 1. 创建创世区块
	genesisBlock := &Block{
		Version:    CURRENT_BLOCK_VERSION,
		PrevHash:   []byte{},
		MerkleRoot: []byte{},
		Timestamp:  uint64(time.Now().Unix()),
		Difficulty: BASE_BLOCK_DIFFCULTY,
		// 其他字段
	}
	pow := new(POW)
	pow.GenerateBlock(genesisBlock)
	// 2. 序列化
	blockData, _ := json.Marshal(genesisBlock)
	// 2. 创建保存目录
	err := os.MkdirAll(filepath.Dir(genesisFile), 0700)
	if err != nil {
		return nil
	}

	// 3. 存储到文件
	err = os.WriteFile(genesisFile, blockData, 0644)
	if err != nil {
		// 错误处理
		return nil
	}
	// 4. 读取并反序列化
	data, _ := os.ReadFile(genesisFile)
	var savedBlock Block
	json.Unmarshal(data, &savedBlock)
	// 5. 验证
	if reflect.DeepEqual(genesisBlock, savedBlock) {
		fmt.Println("创世区块保存成功!")
	}
	return genesisBlock
}
