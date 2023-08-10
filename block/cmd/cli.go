package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/Alan-333333/simple-blockchain/block"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  blockchainCLI [command] [flags]")
	fmt.Println("Commands:")
	fmt.Println("  print - Prints the blockchain")
	fmt.Println("  getblock [hash] - Prints a block")
}

func printBlockchain(chain *block.Blockchain) {
	blocks := chain.Blocks()
	for _, block := range blocks {
		printBlock(block)
	}
}

func printBlock(block *block.Block) {
	fmt.Printf("Hash: %x\n", block.Hash)
	// 打印区块其他信息
	// ...
}

func main() {
	bc := block.NewBlockchain()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "print":
		printBlockchain(bc)

	case "getblock":
		if len(os.Args) < 3 {
			printUsage()
			os.Exit(1)
		}
		hash := os.Args[2]
		block := bc.GetBlock([]byte(hash))
		printBlock(block)

	case "createGenesisBlock":
		// 1. 创建创世区块
		genesisBlock := &block.Block{
			Version:    block.CURRENT_BLOCK_VERSION,
			PrevHash:   []byte{},
			MerkleRoot: []byte{},
			Timestamp:  uint64(time.Now().Unix()),
			// 其他字段
		}
		// 2. 序列化
		blockData, _ := json.Marshal(genesisBlock)
		// 3. 存储到文件
		err := os.WriteFile("genesis.blk", blockData, 0644)
		if err != nil {
			// 错误处理
		}
		// 4. 读取并反序列化
		data, _ := os.ReadFile("genesis.blk")
		var savedBlock block.Block
		json.Unmarshal(data, &savedBlock)
		// 5. 验证
		if reflect.DeepEqual(genesisBlock, savedBlock) {
			fmt.Println("创世区块保存成功!")
		}
		// 6. 添加到区块链
		bc.AddBlock(genesisBlock)
	default:
		printUsage()
		os.Exit(1)
	}
}
