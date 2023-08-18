package main

import (
	"time"

	blockchain "github.com/Alan-333333/simple-blockchain/block/chain"
	"github.com/Alan-333333/simple-blockchain/network/p2p"
	"github.com/Alan-333333/simple-blockchain/transaction"
	"github.com/Alan-333333/simple-blockchain/wallet"
)

var txPool *transaction.TxPool
var bc *blockchain.Blockchain

func main() {

	pow := &blockchain.POW{}
	blockchain.NewBlockchain(pow)
	txPool = transaction.NewTxPool()

	node := p2p.NewNode("127.0.0.1", 3000)
	go node.Listen()

	node.Connect("127.0.0.1", 4000)

	time.Sleep(1 * time.Second)

	// 1. 创建钱包
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	// 2. 构造交易
	tx := transaction.NewTransaction(walletA.GetAddress(), walletB.GetAddress(), 10)

	// 3. 签名
	tx.Sign(walletA.PrivateKey)
	// 4. 广播交易
	node.BroadcastTx(tx)

	time.Sleep(1 * time.Second)
	// 构造交易

	// 5.创建一个区块
	block := blockchain.CreateGenesisBlock()
	block.Transactions = append(block.Transactions, tx)

	//6. 广播区块
	node.BroadcastBlock(block)

	time.Sleep(1 * time.Second)
}
