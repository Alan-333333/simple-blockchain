package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	blockchain "github.com/Alan-333333/simple-blockchain/block/chain"
	"github.com/Alan-333333/simple-blockchain/network/p2p"
	"github.com/Alan-333333/simple-blockchain/transaction"
	"github.com/Alan-333333/simple-blockchain/wallet"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  [command] [flags]")
	fmt.Println("Commands:")
	fmt.Println("  printBlockChain - Prints the blockchain")
	fmt.Println("  getBlock [hash] - Prints a block")
	fmt.Println("  createGenesisBlock - Create a genesis block ")
	fmt.Println("  createwallet - Create a new wallet ")
	fmt.Println("  getwalletbalance [address] - Get wallet balance by address")
	fmt.Println("  send -fromAddress  [address] -toAddress [address] -amount [amount] - Send blockchain transaction")
	fmt.Println("  nodeConnect [IP] [Port]- Connect node")
	fmt.Println("  addwalletBalance [address] [amount] - Add balance fo address")
	fmt.Println()
}

func main() {

	// 区块和交易
	pow := &blockchain.POW{}
	bc := blockchain.NewBlockchain(pow)
	txPool := transaction.NewTxPool()
	bc.Save()

	go bc.Mine(txPool)

	port := 3000
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1])
	}

	// p2p
	node := p2p.NewNode("127.0.0.1", port)
	go node.Listen()

	go startCLI(bc, txPool, node)

	printUsage()

	// 阻塞主线程
	forever := make(chan bool)
	<-forever

}

func printBlockchain(chain *blockchain.Blockchain) {
	blocks := chain.Blocks()
	for _, block := range blocks {
		printBlock(block)
	}
}

func printBlock(block *blockchain.Block) {
	fmt.Printf("Hash: %x\n", block.Hash)
	// 打印区块其他信息
	// ...
}

func startCLI(bc *blockchain.Blockchain, txPool *transaction.TxPool, node *p2p.Node) {
	for {
		// 读取用户输入
		cmd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		args := strings.Split(cmd, " ")

		switch args[0] {
		case "printBlockChain":
			printBlockchain(bc)
		case "getblock":
			hash := args[1]
			block := bc.GetBlock([]byte(hash))
			printBlock(block)

		case "createGenesisBlock":
			// 1. 创建创世区块
			genesisBlock := blockchain.CreateGenesisBlock()
			// 6. 添加到区块链
			err := bc.AddBlock(genesisBlock)
			if err != nil {
				fmt.Println(err)
			}
			bc.Save()
			node.BroadcastBlock(genesisBlock)

			fmt.Println("success")

		case "printNodePeers":
			fmt.Println("peers", node.Server.Peers)

		case "createwallet":
			wallet := wallet.NewWallet()
			wallet.Save()
			fmt.Println("peers", node.Server.Peers)
			node.BroadcastWallet(wallet)
			fmt.Println("success wallet address:", wallet.Address)

		case "nodeConnect":
			ip := args[1]
			portStr := args[2]

			port, _ := strconv.Atoi(portStr)
			go node.Connect(ip, port)

			fmt.Println("success")

		case "getwalletbalance":
			address := args[1]
			balance := wallet.GetAddressBalance(address)
			fmt.Println("success wallet balance:", balance)

		case "addwalletBalance":
			address := args[1]

			amountStr := args[2]
			amountFloat, _ := strconv.ParseFloat(amountStr, 32)
			wallet := wallet.GetwalletByAddress(address)
			wallet.Balance += amountFloat

			wallet.Save()
			node.BroadcastWallet(wallet)
			fmt.Println("success")
		case "send":
			// 获取fromAddress参数值
			fromAddress := args[2]

			// 获取toAddress参数值
			toAddress := args[4]

			// 获取amount参数值
			amountStr := args[6]
			amountFloat, _ := strconv.ParseFloat(amountStr, 32)

			walletFrom := wallet.GetwalletByAddress(fromAddress)
			walletTo := wallet.GetwalletByAddress(toAddress)

			// 构造交易
			tx := transaction.NewTransaction(walletFrom.Address, walletTo.Address, float32(amountFloat))
			// 3. 签名
			tx.Sign(walletFrom.PrivateKey)

			txPool.AddTx(tx)

			node.BroadcastTx(tx)

			// 6. 更新余额
			walletFrom.Balance -= amountFloat
			walletTo.Balance += amountFloat

			walletFrom.Save()
			walletTo.Save()

			node.BroadcastWallet(walletFrom)
			node.BroadcastWallet(walletTo)

			fmt.Println("success")

		default:
			printUsage()
		}
	}
}
