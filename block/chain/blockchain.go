package blockchain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Alan-333333/simple-blockchain/transaction"
	"github.com/Alan-333333/simple-blockchain/wallet"
)

// 定义保存的文件名
const chainFile = "./dat/blockchain/chain.dat"

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
func (bc *Blockchain) GetBlocks() []*Block {
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
		txs := pool.PopTransactions(1)
		if len(txs) == 0 {
			continue
		}
		// 2. 创建新区块
		block := CreateBlock(bc, txs)
		if block == nil {
			continue
		}
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

		pool.RemoveTransactions(txs)

		bc.Save()
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
	if prevBlock != nil {
		block := NewBlock(prevBlock.Hash, prevBlock.Difficulty)
		return block
	}
	return nil
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

	// 1. 序列化区块链
	rawData, err := serialize(bc)
	if err != nil {
		return err
	}

	// 2. 创建保存目录
	err = os.MkdirAll(filepath.Dir(chainFile), 0700)
	if err != nil {
		return err
	}

	// 3. 保存数据和元数据
	err = os.WriteFile(chainFile, rawData, 0644)
	if err != nil {
		return err
	}

	// 4. 保存元数据
	meta := bc.GetMetadata()
	err = meta.Save()

	return err
}

// Load 加载区块链
func LoadBlockchain() (*Blockchain, error) {

	// 1. 读取序列化数据
	raw, err := os.ReadFile(chainFile)
	if err != nil {
		return nil, err
	}
	// 2. 反序列化
	bc, _ := deserialize(raw)

	meta := bc.GetMetadata()
	// 3. 校验元数据的完整性
	err = bc.VerifyMetadata(meta)

	if err != nil {
		return nil, err
	}

	return bc, nil

}

// GetMetadata 获取元数据
func (bc *Blockchain) GetMetadata() *BlockchainMetadata {

	meta := &BlockchainMetadata{}

	if len(bc.blocks) > 0 {
		lastBlock := bc.blocks[len(bc.blocks)-1]
		meta.LastBlockHash = lastBlock.Hash
		meta.BlockCount = len(bc.blocks)
	}

	return meta
}

// VerifyMetadata 验证元数据
func (bc *Blockchain) VerifyMetadata(meta *BlockchainMetadata) error {

	// 检查最后一个区块的hash
	lastBlockHash := bc.blocks[len(bc.blocks)-1].Hash

	if !bytes.Equal(meta.LastBlockHash, lastBlockHash) {
		return errors.New("last block hash mismatch")
	}

	// 检查区块数量
	if meta.BlockCount != len(bc.blocks) {
		return errors.New("block count mismatch")
	}

	return nil

}

// serialize 序列化区块链
func serialize(bc *Blockchain) ([]byte, error) {
	return json.Marshal(&struct {
		Blocks []*Block `json:"blocks"`
		Miner  *Miner   `json:"miner"`
	}{
		Blocks: bc.blocks,
		Miner:  bc.miner,
	})
}

func deserialize(data []byte) (*Blockchain, error) {
	var raw struct {
		Blocks []*Block
		Miner  *Miner
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	bc := &Blockchain{
		blocks: raw.Blocks,
		miner:  raw.Miner,
	}

	return bc, nil
}
