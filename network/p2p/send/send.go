package main

import (
	"time"

	"github.com/Alan-333333/simple-blockchain/network/p2p"
)

func main() {
	node := p2p.NewNode("127.0.0.1", 3000)
	// go node.Listen()

	node.Connect("127.0.0.1", 4000)

	time.Sleep(10 * time.Second)
}
