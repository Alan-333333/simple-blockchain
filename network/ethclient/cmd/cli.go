package main

import (
	"fmt"

	ethclient "github.com/Alan-333333/simple-blockchain/network/ethclient/client"
)

const URL = "https://eth-mainnet.g.alchemy.com/v2/lomkr9JTMgtizYfI_ZfvkVzxs_L_EvlZ"

func main() {

	client, err := ethclient.NewEthClient(URL)
	if err != nil {
		panic(err)
	}

	lastBlock, err := client.GetLatestBlock()
	if err != nil {
		panic(err)
	}
	fmt.Println(lastBlock)

}
