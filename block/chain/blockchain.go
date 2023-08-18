package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Alan-333333/simple-blockchain/transaction"
	"github.com/Alan-333333/simple-blockchain/wallet"
	"github.com/google/uuid"
)

var blockchainInstance *Blockchain

type Miner struct {
	wallet *wallet.Wallet // 矿工钱包
}

type Blockchain struct {
	blocks    []*Block
	miner     *Miner
	consensus Consensus
}

// 创建区块链
func NewBlockchain(consensus Consensus) *Blockchain {
	if blockchainInstance != nil {
		return blockchainInstance
	}
	// ...初始化
	blockchainInstance = &Blockchain{
		blocks:    []*Block{},
		consensus: consensus,
	}

	return blockchainInstance
}

func GetBlockchain() *Blockchain {
	return blockchainInstance
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
	if !bc.IsValidBlock(block) {
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
// IsValidBlock 验证区块是否合法
func (bc *Blockchain) IsValidBlock(block *Block) bool {
	return bc.consensus.VerifyBlock(block)
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
		block := CreateBlock(bc, txs)

		// 3. 挖矿
		// doPoW(block)
		bc.consensus.GenerateBlock(block)

		// 填充交易
		block.Transactions = txs

		// 4. 添加区块
		err := bc.AddBlock(block)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func NewEmptyBlock(bc *Blockchain) *Block {
	prevBlock := bc.GetLastBlock()
	block := NewBlock(prevBlock.Hash, prevBlock.Difficulty)
	return block
}

// 创建新区块
func CreateBlock(bc *Blockchain, txs []*transaction.Transaction) *Block {

	prevBlock := bc.GetLastBlock()
	block := NewBlock(prevBlock.Hash, prevBlock.Difficulty)

	return block

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
