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
	fmt.Println("  printBlockChain - Print all blocks in the blockchain")
	fmt.Println("  printBlock [hash] - Print a specific block")
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
	fmt.Printf("Block Version: %v\n", block.Version)

	fmt.Printf("Block PrevHash: %v\n", block.PrevHash)

	fmt.Printf("Block Difficulty: %v\n", block.Difficulty)

	fmt.Printf("Block Transactions: %v\n", block.Transactions)
	// ...

}

// startCLI starts the command line interface
func startCLI(bc *blockchain.Blockchain, txPool *transaction.TxPool, node *p2p.Node) {
	for {
		// Parse input
		args := parseInput()

		switch args.command {

		// Print block chain
		case "printBlockChain":
			printBlockchain(bc)

			// Print block by hash
		case "printBlock":
			// parse hash by args.params
			hash := args.params[0]
			// get block by hash
			block := bc.GetBlock(hash)

			printBlock(block)

			// Create genesis block
		case "createGenesisBlock":
			// Create genesis block
			genesisBlock := blockchain.CreateGenesisBlock()

			// Add block to blockchain
			err := bc.AddBlock(genesisBlock)
			if err != nil {
				fmt.Println(err)
			}

			// Save blockchain
			bc.Save()

			// Broadcast block
			node.BroadcastBlock(genesisBlock)

			// Print success message
			printSuccess()

			// Print node peers
		case "printNodePeers":
			fmt.Println("Peers:", node.Server.Peers)

			// Create new wallet
		case "createWallet":
			// Create new wallet
			wallet := wallet.NewWallet()

			// Save updated wallet
			wallet.Save()

			// Broadcast updated wallet
			node.BroadcastWallet(wallet)

			// Print success message
			fmt.Println("success wallet address:", wallet.Address)
		case "connectNode":
			ip := args.params[0]
			portStr := args.params[1]

			port, _ := strconv.Atoi(portStr)
			go node.Connect(ip, port)

			fmt.Println("success")

			// Get wallet balance
		case "getWalletBalance":
			// Parse wallet address
			address := args.params[0]
			// Get balance
			balance := wallet.GetAddressBalance(address)
			// Print balance
			fmt.Println("success wallet balance:", balance)

			// Add balance to wallet
		case "addWalletBalance":
			address := args.params[1]
			amount := args.params[2]
			amountFloat, _ := strconv.ParseFloat(amount, 32)

			// Get wallet
			wallet := wallet.GetwalletByAddress(address)
			// Update balance
			balance := wallet.Balance + amountFloat
			wallet.UpdateWalletBalance(balance)

			// Save updated wallet
			wallet.Save()

			// Broadcast updated wallet
			node.BroadcastWallet(wallet)

			// Print success message
			printSuccess()

			// Send transaction
		case "sendTransaction":
			// Parse transaction parameters
			fromAddress := parseFromAddress(args)
			toAddress := parseToAddress(args)
			amount := parseAmount(args)

			// Get sender and recipient wallets
			senderWallet := wallet.GetwalletByAddress(fromAddress)
			recipientWallet := wallet.GetwalletByAddress(toAddress)

			// Create new transaction
			tx := transaction.NewTransaction(senderWallet.Address, recipientWallet.Address, float32(amount))

			// Sign the transaction
			tx.Sign(senderWallet.PrivateKey)

			// Add transaction to transaction pool
			txPool.AddTx(tx)

			// Broadcast transaction to network
			node.BroadcastTx(tx)

			// Update wallet balances
			senderWallet.Balance -= amount
			recipientWallet.Balance += amount

			// Save updated wallets
			senderWallet.Save()
			recipientWallet.Save()

			// Broadcast updated wallets to network
			node.BroadcastWallet(senderWallet)
			node.BroadcastWallet(recipientWallet)

			// Print success message
			printSuccess()

		default:
			printUsage()
		}
	}
}

// print success
func printSuccess() {
	fmt.Println("success")
}

// parse from address from input
func parseFromAddress(args Input) string {
	return args.params[1]
}

// parse to address from input
func parseToAddress(args Input) string {
	return args.params[3]
}

// parse amount from input
func parseAmount(args Input) float64 {
	amountStr := args.params[4]
	amountFloat, _ := strconv.ParseFloat(amountStr, 32)
	return amountFloat
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
