package main

import (
	"fmt"
	"log"
	"math/big"
	"os"

	ethclient "github.com/Alan-333333/simple-blockchain/network/ethclient/client"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const URL = "https://eth-mainnet.g.alchemy.com/v2/lomkr9JTMgtizYfI_ZfvkVzxs_L_EvlZ"

func main() {
	client, err := ethclient.NewEthClient(URL)
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Println("请选择操作:")
		fmt.Println("1. 获取最新区块号")
		fmt.Println("2. 获取账户余额")
		fmt.Println("3. 发送交易")
		fmt.Println("4. 获取交易详情")
		fmt.Println("5. 获取交易收据")
		fmt.Println("6. 过滤日志")
		fmt.Println("7. 退出")

		var input string
		fmt.Scanln(&input)

		switch input {
		case "1":
			blockNum, err := client.GetLatestBlockNumber()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("最新区块号:", blockNum)

		case "2":
			var account string
			fmt.Print("输入账户地址:")
			fmt.Scanln(&account)

			address := common.HexToAddress(account)

			balance, err := client.GetBalance(address)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("账户余额:", balance)

		case "3":

			// 1. 提示用户构造交易
			fmt.Print("输入发送者地址:")
			var fromStr string
			fmt.Scanln(&fromStr)
			from := common.HexToAddress(fromStr)
			// 获取from的字节数组表示
			fromBytes := from.Bytes()

			fmt.Print("输入接收者地址:")
			var toStr string
			fmt.Scanln(&toStr)
			to := common.HexToAddress(toStr)

			fmt.Print("输入发送金额:")
			var amountStr string
			fmt.Scanln(&amountStr)
			amount, _ := new(big.Int).SetString(amountStr, 10)

			fmt.Print("输入nonce:")
			var nonceStr string
			fmt.Scanln(&nonceStr)
			nonce, _ := new(big.Int).SetString(nonceStr, 10)

			fmt.Print("输入gas price:")
			var gasPriceStr string
			fmt.Scanln(&gasPriceStr)
			gasPrice, _ := new(big.Int).SetString(gasPriceStr, 10)

			// 2. 构造交易
			nonceUint64 := nonce.Uint64()
			tx := types.NewTransaction(nonceUint64, to, amount, 21000, gasPrice, fromBytes)

			// 3. 发送交易

			err = client.SendTransaction(tx)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("交易已发送")

		case "4":
			fmt.Print("输入交易哈希:")
			var txHashStr string
			fmt.Scanln(&txHashStr)

			txHash := common.HexToHash(txHashStr)

			tx, pending, err := client.GetTransactionByHash(txHash)
			if err != nil {
				log.Fatal(err)
			}

			if pending {
				fmt.Println("交易正在处理")
			} else {
				fmt.Println("交易详情:", tx)
			}

		case "5":
			fmt.Print("输入交易哈希:")
			var txHashStr string
			fmt.Scanln(&txHashStr)
			txHash := common.HexToHash(txHashStr)

			receipt, _ := client.GetTransactionReceipt(txHash)
			fmt.Println("交易收据:", receipt)
		case "6":
			fmt.Print("输入监听地址:")
			var addressStr string
			fmt.Scanln(&addressStr)
			address := common.HexToAddress(addressStr)

			query := ethereum.FilterQuery{
				Addresses: []common.Address{address},
			}

			logs, _ := client.FilterLogs(query)
			fmt.Println("相关日志:", logs)

		case "7":
			os.Exit(0)
		}
	}

}
