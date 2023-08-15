package main

import "github.com/Alan-333333/simple-blockchain/network/p2p"

func main() {
	node := p2p.NewNode("127.0.0.1", 4000)
	node.Listen()
}
