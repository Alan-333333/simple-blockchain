package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

const CURRENT_BLOCK_VERSION = 1

type Blockchain struct {
	blocks []*Block
}

// 创建区块链
func NewBlockchain() *Blockchain {
	return &Blockchain{
		blocks: []*Block{},
	}
}

// 返回区块链中所有的区块
func (bc *Blockchain) Blocks() []*Block {
	blocks := make([]*Block, len(bc.blocks))
	copy(blocks, bc.blocks)

	return blocks
}

// 添加新区块
// AddBlock 向区块链中添加新区块
func (bc *Blockchain) AddBlock(block *Block) error {

	// 验证新区块
	if !isValidBlock(block) {
		return fmt.Errorf("block is not a valid block")
	}

	// 获取上一个区块并设置上一个hash
	prevBlock := bc.blocks[len(bc.blocks)-1]
	block.PrevHash = prevBlock.Hash

	// 添加区块
	bc.blocks = append(bc.blocks, block)

	return nil
}

// 获取最后一个区块
func (bc *Blockchain) GetLastBlock() *Block {
	blocks := bc.blocks

	// 区块链为空
	if len(blocks) == 0 {
		return nil
	}

	// 取最后一个区块
	return blocks[len(blocks)-1]
}

// 判断区块是否valid
// isValidBlock 验证区块是否合法
func isValidBlock(block *Block) bool {

	// 基本参数校验
	if block.Version != CURRENT_BLOCK_VERSION {
		return false
	}

	if block.Timestamp > uint64(time.Now().Unix()) {
		return false
	}

	// 验证交易的合法性
	for _, tx := range block.Transactions {
		if !IsValidTransaction(tx) {
			return false
		}
	}

	// 验证区块Hash
	blockHash := calcBlockHash(block)
	if !bytes.Equal(blockHash, block.Hash) {
		return false
	}

	// 其他规则校验
	// ...

	return true
}

// IsValidTransaction 交易合法性校验
func IsValidTransaction(tx *Transaction) bool {

	// 交易基本校验,Sender, Receiver, Value
	if tx.From == nil || tx.To == nil {
		return false
	}

	if tx.Value <= 0 {
		return false
	}

	// 其他规则校验
	// ...

	return true

}

// 计算区块hash
func calcBlockHash(block *Block) []byte {

	// 1. 序列化区块数据
	blockData, err := json.Marshal(block)
	if err != nil {
		panic("序列化区块数据失败")
	}

	// 2. 生成SHA256哈希
	blockHash := sha256.Sum256(blockData)

	// 3. 返回字节数组
	return blockHash[:]
}

// 实现共识算法
func RunConsensus(newBlock *Block) error {

	// 本地挖矿
	for {

		// 随机生成nonce值
		nonce := rand.Int63()
		newBlock.Nonce = make([]byte, 8)
		binary.BigEndian.PutUint64(newBlock.Nonce, uint64(nonce))

		// 计算新的hash
		newHash := calcBlockHash(newBlock)

		// 检查hash是否满足难度要求
		if meetsDifficulty(newHash, newBlock.Difficulty) {
			return nil
		}
	}

}

func meetsDifficulty(hash []byte, difficulty uint64) bool {

	// 简单检查hash的前n位是否为0
	// n由difficulty决定
	prefix := bytes.Repeat([]byte{0}, int(difficulty/8))
	return bytes.HasPrefix(hash, prefix)

}

// 在区块链中根据hash获取区块
func (bc *Blockchain) GetBlock(blockHash []byte) *Block {

	for _, block := range bc.blocks {
		if bytes.Equal(block.Hash, blockHash) {
			return block
		}
	}

	return nil
}

// 在区块链中根据高度获取区块
func (bc *Blockchain) GetBlockByHeight(height int) *Block {

	if height > len(bc.blocks) {
		return nil
	}

	return bc.blocks[height]
}
