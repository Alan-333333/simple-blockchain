package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/Alan-333333/simple-blockchain/utils"
)

type Transaction struct {
	Sender    string
	Recipient string
	Value     float32
	Signature string

	// 其他字段
}

// NewTransaction 创建新交易
func NewTransaction(sender string, recipient string, value float32) *Transaction {

	tx := &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Value:     value,
	}

	return tx

}

// Sign 交易签名
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) error {
	// 签名算法
	signature := signECDSA(privateKey, tx) // 使用私钥签名

	tx.Signature = signature // 添加签名到交易

	return nil
}

// IsValid 验证交易签名
func (tx *Transaction) IsValid() bool {
	// 解析公钥
	fmt.Println(tx.Sender)
	pubKey, err := utils.AddrToPubKey(tx.Sender)
	fmt.Println(pubKey)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// 签名验证算法
	return verifyECDSA(pubKey, tx, tx.Signature) // 使用公钥验证
}

// 对交易签名
func signECDSA(privateKey *ecdsa.PrivateKey, tx *Transaction) string {

	// 序列化交易
	txBytes := serialize(tx)

	// 计算交易的哈希
	txHash := sha256.Sum256(txBytes)

	// 签名交易哈希
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, txHash[:])
	if err != nil {
		return ""
	}

	// 将签名转换为16进制字符串
	sigStr := fmt.Sprintf("%064x%064x", r, s)

	return sigStr
}

// 验证签名
func verifyECDSA(publicKey *ecdsa.PublicKey, tx *Transaction, sig string) bool {

	// 解析签名字符串
	r, s := utils.ParseSig(sig)

	// 序列化交易
	txBytes := serialize(tx)

	// 计算交易哈希
	txHash := sha256.Sum256(txBytes)

	// 使用公钥验证签名
	return ecdsa.Verify(publicKey, txHash[:], r, s)

}

func serialize(tx *Transaction) []byte {
	// 序列化交易为字节数组
	return []byte(fmt.Sprintf("%v%v", tx))
}
