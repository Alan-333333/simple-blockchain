package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Alan-333333/simple-blockchain/transaction"
	"github.com/Alan-333333/simple-blockchain/wallet"
	"github.com/google/uuid"
)

type Miner struct {
	wallet *wallet.Wallet // 矿工钱包
}

type Blockchain struct {
	blocks []*Block
	miner  *Miner
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
	if len(bc.blocks) != 0 {
		prevBlock := bc.blocks[len(bc.blocks)-1]
		block.PrevHash = prevBlock.Hash
	}

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
	// blockHash := calcBlockHash(block)
	// if !bytes.Equal(blockHash, block.Hash) {
	// 	return false
	// }

	// 其他规则校验
	// ...

	return true
}

// IsValidTransaction 交易合法性校验
func IsValidTransaction(tx *transaction.Transaction) bool {

	// 交易基本校验,Sender, Receiver, Value
	if tx.Sender == "" || tx.Recipient == "" {
		return false
	}

	if tx.Value <= 0 {
		return false
	}

	// 其他规则校验
	// ...

	return true

}

// 挖矿
func (bc *Blockchain) Mine(pool *transaction.TxPool) {
	for {
		// 1.获取新的交易
		txs := pool.GetTxs()
		// 2. 创建新区块
		block := createBlock(bc, txs)

		// 3. 挖矿
		doPoW(block)

		// 4. 添加区块
		err := bc.AddBlock(block)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// 创建新区块
func createBlock(bc *Blockchain, txs []*transaction.Transaction) *Block {

	prevBlock := bc.GetLastBlock()
	block := NewBlock(prevBlock.Hash, prevBlock.Difficulty)
	// 填充交易
	block.Transactions = txs

	return block

}

// 实现共识算法
func doPoW(newBlock *Block) {

	// 本地挖矿

	// 随机生成nonce值
	nonce := rand.Int63()
	newBlock.Nonce = make([]byte, 8)
	binary.BigEndian.PutUint64(newBlock.Nonce, uint64(nonce))

	// 计算新的hash
	newHash := calcBlockHash(newBlock)

	// 检查hash是否满足难度要求
	if meetsDifficulty(newHash, newBlock.Difficulty) {
		return
	}

	newBlock.Hash = newHash
	// 难度递增
	newBlock.Difficulty = newBlock.Difficulty + 1

}

// 计算区块hash
func calcBlockHash(block *Block) []byte {

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

func meetsDifficulty(hash []byte, difficulty uint64) bool {

	// 简单检查hash的前n位是否为0
	// n由difficulty决定
	prefix := bytes.Repeat([]byte{0}, int(difficulty/8))
	return bytes.HasPrefix(hash, prefix)

}

// func getDifficultyPrefix(difficulty uint64) []byte {

// 	// 难度控制0的位数
// 	zeroNum := difficulty / 8

// 	// 生成字节数组
// 	prefix := make([]byte, zeroNum)

// 	// 填充0
// 	for i := 0; i < zeroNum; i++ {
// 		prefix[i] = byte(0)
// 	}

// 	return prefix

// }

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

// Save 将区块链序列化为文件
func (bc *Blockchain) Save() error {

	// 1. 将区块链序列化为字节数组
	rawData, err := serialize(bc)
	if err != nil {
		return err
	}

	// 2. 将数据写入文件
	uuid := uuid.New()
	dbPath := "chain-" + uuid.String() + ".dat"
	err = os.WriteFile(dbPath, rawData, 0644)
	if err != nil {
		return err
	}

	return nil

}

// serialize 序列化区块链
func serialize(bc *Blockchain) ([]byte, error) {
	rawData, err := json.Marshal(bc)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

// LoadBlockchain 加载持久化的区块链
func LoadBlockchain(file string) *Blockchain {

	rawData, _ := os.ReadFile(file)

	bc := deserialize(rawData)

	return bc
}

func deserialize(data []byte) *Blockchain {
	var bc Blockchain
	err := json.Unmarshal(data, &bc)
	if err != nil {
		panic(err)
	}
	return &bc
}
