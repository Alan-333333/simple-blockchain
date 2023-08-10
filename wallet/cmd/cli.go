package main

import (
	"fmt"
	"log"

	blockchain "github.com/Alan-333333/simple-blockchain/block/chain"
	"github.com/Alan-333333/simple-blockchain/transaction"
	"github.com/Alan-333333/simple-blockchain/wallet"
)

func main() {

	simulate()

}

func simulate() {

	// 1. 创建钱包
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	// 2. 构造交易
	tx := transaction.NewTransaction(walletA.GetAddress(), walletB.GetAddress(), 10)

	// 3. 签名
	tx.Sign(walletA.PrivateKey)

	// 4. 验证
	if !tx.IsValid() {
		log.Fatal("Invalid transaction")
	}

	pool := transaction.NewTxPool()
	pool.AddTx(tx)

	newBlock := blockchain.NewBlock()

	bc := blockchain.NewBlockchain()
	bc.AddBlock(newBlock)

	// 5. 模拟执行
	fmt.Println("Transfer 10 coins from", walletA.GetAddress(), "to", walletB.GetAddress())

	// 6. 更新余额
	walletA.Balance -= 10
	walletB.Balance += 10

	walletA.Save()
	walletB.Save()

	bc.Save()
}
