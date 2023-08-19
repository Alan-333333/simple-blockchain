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

type Input struct {
	command string
	params  []string
}

// printUsage prints help text explaining the program usage and commands
func printUsage() {

	fmt.Println("Usage:")
	fmt.Println("  [command] [flags]")

	// Print blockchain related commands
	fmt.Println("Blockchain Commands:")
	fmt.Println("  printBlockchain - Print all blocks in the blockchain")
	fmt.Println("  getBlock [hash] - Print a specific block")
	fmt.Println("  createGenesisBlock - Create the genesis block")

	// Print wallet related commands
	fmt.Println("Wallet Commands:")
	fmt.Println("  createWallet - Create a new wallet")
	fmt.Println("  getWalletBalance [address] - Get balance for a wallet")
	fmt.Println("  addWalletBalance [address] [amount] - Add balance to a wallet")

	// Print transaction related commands
	fmt.Println("Transaction Commands:")
	fmt.Println("  sendTransaction -from [address] -to [address] -amount [amount] - Send a transaction")

	// Print node related commands
	fmt.Println("Node Commands:")
	fmt.Println("  connectNode [IP] [port] - Connect to a node")

	// Print node related commands
	fmt.Println()

}

// main is the entry point of the program
func main() {

	// Initialize blockchain
	pow := &blockchain.POW{}
	bc := blockchain.NewBlockchain(pow)
	txPool := transaction.NewTxPool()
	bc.Save()

	// Start mining
	go bc.Mine(txPool)

	// Parse command line args
	port := 3000
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1])
	}

	// Start P2P node
	node := p2p.NewNode("127.0.0.1", port)
	go node.Listen()

	// Start CLI
	go startCLI(bc, txPool, node)

	// Print usage
	printUsage()

	// Block main thread
	forever := make(chan bool)
	<-forever
}

// printBlockchain prints all blocks in the blockchain
func printBlockchain(chain *blockchain.Blockchain) {

	// Get all blocks
	blocks := chain.GetBlocks()

	// Print each block
	for _, block := range blocks {
		printBlock(block)
	}

}

// printBlock prints a single block's data
func printBlock(block *blockchain.Block) {

	// Print block hash
	fmt.Printf("Block Hash: %x\n", block.Hash)

	// Print other block data
	// ...

}

// startCLI starts the command line interface
func startCLI(bc *blockchain.Blockchain, txPool *transaction.TxPool, node *p2p.Node) {
	for {
		// Parse input
		args := parseInput()

		switch args.command {

		case "printBlockChain":
			printBlockchain(bc)

		case "getblock":
			hash := args.params[0]
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
			ip := args.params[0]
			portStr := args.params[1]

			port, _ := strconv.Atoi(portStr)
			go node.Connect(ip, port)

			fmt.Println("success")

		case "getwalletbalance":
			address := args.params[0]
			balance := wallet.GetAddressBalance(address)
			fmt.Println("success wallet balance:", balance)

		case "addwalletBalance":
			address := args.params[1]

			amountStr := args.params[1]
			amountFloat, _ := strconv.ParseFloat(amountStr, 32)
			wallet := wallet.GetwalletByAddress(address)
			wallet.Balance += amountFloat

			wallet.Save()
			node.BroadcastWallet(wallet)
			fmt.Println("success")

		case "send":
			// 获取fromAddress参数值
			fromAddress := args.params[1]

			// 获取toAddress参数值
			toAddress := args.params[3]

			// 获取amount参数值
			amountStr := args.params[4]
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

// parseInput parses user input into command and parameters
func parseInput() Input {

	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	cmd := strings.TrimSpace(input)

	params := strings.Split(cmd, " ")

	return Input{
		command: params[0],
		params:  params[1:],
	}
}
