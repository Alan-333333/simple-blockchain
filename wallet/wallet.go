package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Alan-333333/simple-blockchain/utils"
)

const WALLET_PATE = "wallet.dat"

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
	Address    string  // 地址就是公钥的Hash
	Balance    float64 // 新增余额字段
}

func NewWallet() *Wallet {

	// 1. 生成私钥
	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	// 2. 从私钥计算公钥
	pubKey := privKey.PublicKey

	// 3. 生成地址(公钥hash)
	address := utils.PubKeyToAddr(&pubKey)

	return &Wallet{
		PrivateKey: privKey,
		PublicKey:  &pubKey,
		Address:    address,
		Balance:    0, // 初始化余额
	}
}

func (w *Wallet) GetAddress() string {
	return w.Address
}

func (wallet *Wallet) Save() error {

	// 1. 序列化钱包数据
	jsonData := serialize(wallet)

	// 2. 将数据写入文件

	path := fmt.Sprintf("%s.wallet", wallet.Address)

	err := os.WriteFile(path, []byte(jsonData), 0600)

	return err

}

func serialize(wallet *Wallet) string {
	json, _ := json.Marshal(wallet)
	return string(json)
}

func deserialize(data string) *Wallet {
	var wallet Wallet
	json.Unmarshal([]byte(data), &wallet)
	return &wallet
}
