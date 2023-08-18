package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"time"
)

type POW struct{}

func (pow *POW) GenerateBlock(newBlock *Block) {

	// 随机生成nonce值
	nonce := rand.Int63()
	newBlock.Nonce = make([]byte, 8)
	binary.BigEndian.PutUint64(newBlock.Nonce, uint64(nonce))

	// 计算新的hash
	newHash := pow.CalcBlockHash(newBlock)

	// 检查hash是否满足难度要求
	if meetsDifficulty(newHash, newBlock.Difficulty) {
		return
	}

	newBlock.Hash = newHash
	// 难度递增
	newBlock.Difficulty = newBlock.Difficulty + 1
}

func (p *POW) VerifyBlock(block *Block) bool {
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
	blockHash := p.CalcBlockHash(block)
	return bytes.Equal(blockHash, block.Hash)
}

func meetsDifficulty(hash []byte, difficulty uint64) bool {

	// 简单检查hash的前n位是否为0
	// n由difficulty决定
	prefix := bytes.Repeat([]byte{0}, int(difficulty/8))
	return bytes.HasPrefix(hash, prefix)

}

// 计算区块hash
func (pow *POW) CalcBlockHash(block *Block) []byte {
	tsBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(tsBytes, block.Timestamp)
	// 拼接区块数据
	data := bytes.Join([][]byte{
		block.PrevHash,
		block.MerkleRoot,
		tsBytes,
	}, []byte{})

	// 拼接nonce
	data = append(data, block.Nonce...)

	// 计算hash
	hash := sha256.Sum256(data)
	// 返回字节数组
	return hash[:]

}
