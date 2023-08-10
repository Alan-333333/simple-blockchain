package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
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

	fmt.Println(pubKey)
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

func (w *Wallet) Sign(data []byte) string {

	// 1. 对数据hash
	hash := sha256.Sum256(data)

	// 2. 签名
	r, s, _ := ecdsa.Sign(rand.Reader, w.PrivateKey, hash[:])

	// 3. 序列化签名
	sig := fmt.Sprintf("%064x%064x", r, s)

	return sig
}

func Verify(pubkey *ecdsa.PublicKey, data, sig string) bool {

	// 1. 解析签名
	r, s := utils.ParseSig(sig)

	// 2. 对数据hash
	hash := sha256.Sum256([]byte(data))

	// 3. 验证签名
	return ecdsa.Verify(pubkey, hash[:], r, s)
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
