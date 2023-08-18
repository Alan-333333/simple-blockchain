package main

import (
	blockchain "github.com/Alan-333333/simple-blockchain/block/chain"
	"github.com/Alan-333333/simple-blockchain/network/p2p"
	"github.com/Alan-333333/simple-blockchain/transaction"
)

func main() {

	// 区块和交易
	pow := &blockchain.POW{}
	bc := blockchain.NewBlockchain(pow)
	txPool := transaction.NewTxPool()

	go bc.Mine(txPool)

	// p2p
	node := p2p.NewNode("127.0.0.1", 3000)
	go node.Listen()

}
