package main

import (
	"fmt"
	"log"
	"time"

	blockchain "github.com/Alan-333333/simple-blockchain/block/chain"
	"github.com/Alan-333333/simple-blockchain/transaction"
	"github.com/Alan-333333/simple-blockchain/wallet"
)

func main() {
	// 1. 创建创世块
	genesisBlock := blockchain.CreateGenesisBlock()

	// 2. 创建区块链并添加创世块
	chain := blockchain.NewBlockchain()
	err := chain.AddBlock(genesisBlock)
	if err != nil {
		log.Fatalln(err)
	}

	// 5. 创建交易池并添加交易
	pool := transaction.NewTxPool()

	go chain.Mine(pool)

	// 3. 创建2个账户
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	// 4. walletA向walletB转账10个币
	tx := transaction.NewTransaction(walletA.Address, walletB.Address, 10)
	tx.Sign(walletA.PrivateKey)

	pool.AddTx(tx)

	time.Sleep(100 * time.Millisecond)
	// 5. 验证
	block := chain.GetLastBlock()

	fmt.Println(block.Transactions[0]) // 打印交易
	fmt.Println(walletB.Balance)       // 显示接收方余额

	// 再创建一笔交易
	tx = transaction.NewTransaction(walletB.Address, walletA.Address, 10)
	tx.Sign(walletB.PrivateKey)

	pool.AddTx(tx)

	// 6. 验证
	time.Sleep(100 * time.Millisecond)
	block = chain.GetLastBlock()

	fmt.Println(block.Transactions[0]) // 打印交易
	fmt.Println(walletB.Balance)       // 显示接收方余额

}
