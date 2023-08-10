package blockchain

import "github.com/Alan-333333/simple-blockchain/tx"

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
	Transactions []*tx.Transaction
}
