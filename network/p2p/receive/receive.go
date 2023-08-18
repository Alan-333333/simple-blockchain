package main

import (
	blockchain "github.com/Alan-333333/simple-blockchain/block/chain"
	"github.com/Alan-333333/simple-blockchain/network/p2p"
	"github.com/Alan-333333/simple-blockchain/transaction"
)

func main() {

	pow := &blockchain.POW{}
	blockchain.NewBlockchain(pow)

	transaction.NewTxPool()

	node := p2p.NewNode("127.0.0.1", 4000)
	go node.Listen()

	for {

	}
}
