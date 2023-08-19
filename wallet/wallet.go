package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

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

	path := fmt.Sprintf("./dat/wallet/%s.wallet", wallet.Address)

	// 创建包含路径的文件夹
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	err := os.WriteFile(path, []byte(jsonData), 0600)

	return err

}

// 查询钱包
func GetwalletByAddress(address string) *Wallet {

	// 1. 打开钱包文件
	walletFile := fmt.Sprintf("./dat/wallet/%s.wallet", address)
	fileData, err := os.ReadFile(walletFile)
	if err != nil {
		return nil
	}

	// 2. 反序列化数据到钱包结构体
	wallet := deserialize(string(fileData))

	// 3. 返回钱包的余额
	return wallet
}

// 查询地址余额
func GetAddressBalance(address string) float64 {

	// 1. 打开钱包文件
	walletFile := fmt.Sprintf("./dat/wallet/%s.wallet", address)
	fileData, err := os.ReadFile(walletFile)
	if err != nil {
		return 0
	}

	// 2. 反序列化数据到钱包结构体
	wallet := deserialize(string(fileData))

	// 3. 返回钱包的余额
	return wallet.Balance
}

func serialize(wallet *Wallet) string {
	json, _ := EncodedWallet(wallet)
	return string(json)
}

func deserialize(data string) *Wallet {

	wallet, _ := DecodeWallet([]byte(data))
	return wallet
}

type NodeWallet struct {
	PublicKey  []byte
	PrivateKey []byte
	Address    string
	Balance    float64
}

func EncodedWallet(wallet *Wallet) ([]byte, error) {

	pubASN1, err := x509.MarshalPKIXPublicKey(wallet.PublicKey)
	if err != nil {
		return []byte{}, nil
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})

	priASN1, err := x509.MarshalECPrivateKey(wallet.PrivateKey)

	if err != nil {
		return []byte{}, nil
	}
	priByte := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: priASN1,
	})

	ewallet := NodeWallet{
		PublicKey:  pubBytes,
		PrivateKey: priByte,
		Address:    wallet.Address,
		Balance:    wallet.Balance,
	}

	return json.Marshal(ewallet)
}

// 解码Wallet
func DecodeWallet(data []byte) (*Wallet, error) {

	var ewallet NodeWallet

	err := json.Unmarshal(data, &ewallet)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(ewallet.PublicKey)
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	block, _ = pem.Decode(ewallet.PrivateKey)
	privKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	wallet := &Wallet{
		PublicKey:  pubKey.(*ecdsa.PublicKey),
		PrivateKey: privKey,
		Address:    ewallet.Address,
		Balance:    ewallet.Balance,
	}

	return wallet, nil
}
