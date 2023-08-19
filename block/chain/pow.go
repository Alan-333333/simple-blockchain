package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

type POW struct{}

func (pow *POW) GenerateBlock(newBlock *Block) {

	// 随机生成nonce值
	nonce := rand.Int63()
	newBlock.Nonce = make([]byte, 8)
	binary.BigEndian.PutUint64(newBlock.Nonce, uint64(nonce))

	for {
		// 在一定范围内调整nonce
		nonce := nextNonce(newBlock.Nonce)
		newBlock.Nonce = nonce

		hash := CalcBlockHash(newBlock)

		if meetsDifficulty(hash, newBlock.Difficulty) {
			newBlock.Hash = hash
			return
		}
		// 难度不能增长太快
		newBlock.Difficulty = adjustDifficulty(newBlock.Difficulty)
	}
}

func nextNonce(nonce []byte) []byte {
	// 将nonce解析为big.Int
	n := new(big.Int).SetBytes(nonce)

	// 在原nonce附近做一个小的调整
	n = n.Add(n, big.NewInt(1))

	// 转换回字节数组
	return n.Bytes()
}

func adjustDifficulty(diff uint64) uint64 {
	// 控制难度增长速度
	return diff + 1
}

func (p *POW) VerifyBlock(block *Block) bool {
	// 基本参数校验
	if block.Version != CURRENT_BLOCK_VERSION {
		fmt.Println("err version")
		return false
	}

	if block.Timestamp > uint64(time.Now().Unix()) {
		fmt.Println("err Timestamp")
		return false
	}

	// 验证交易的合法性
	for _, tx := range block.Transactions {
		if !IsValidTransaction(tx) {
			fmt.Println("err Transactions")
			return false
		}
	}

	// 验证区块Hash
	blockHash := CalcBlockHash(block)
	result := bytes.Equal(blockHash, block.Hash)
	if !result {
		fmt.Println("err Hash,", blockHash)
		fmt.Println("err Hash,", block.Hash)
	}
	return result
}

func meetsDifficulty(hash []byte, difficulty uint64) bool {

	// 简单检查hash的前n位是否为0
	// n由difficulty决定
	prefix := bytes.Repeat([]byte{0}, int(difficulty/8))
	return bytes.HasPrefix(hash, prefix)

}

// 计算区块hash
func CalcBlockHash(block *Block) []byte {
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
