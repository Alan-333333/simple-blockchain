package transaction

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"fmt"
	"math/big"

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
	pubKey, err := utils.AddrToPubKey(tx.Sender)
	if err != nil {
		return false
	}
	// 签名验证算法
	return verifyECDSA(pubKey, tx, tx.Signature) // 使用公钥验证
}

// 对交易签名
func signECDSA(privateKey *ecdsa.PrivateKey, tx *Transaction) string {
	// 序列化交易
	txBytes := serialize(tx)
	// 计算哈希
	txHash := sha256.Sum256(txBytes)

	// 签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, txHash[:])
	if err != nil {
		return ""
	}
	// 序列化签名为字节数组
	sigBytes, err := asn1.Marshal(struct{ R, S *big.Int }{r, s})
	if err != nil {
		return ""
	}

	// 编码为Base64
	return base64.StdEncoding.EncodeToString(sigBytes)
}

// 验证签名
func verifyECDSA(publicKey *ecdsa.PublicKey, tx *Transaction, sigStr string) bool {
	// 1. Base64解码
	sigBytes, err := base64.StdEncoding.DecodeString(sigStr)
	if err != nil {
		return false
	}

	// 2. ASN.1反序列化
	var sig struct{ R, S *big.Int }
	if _, err := asn1.Unmarshal(sigBytes, &sig); err != nil {
		return false
	}

	// 3. 序列化交易
	txBytes := serialize(tx)
	// 4. 计算哈希
	txHash := sha256.Sum256(txBytes)

	// 5. 验证签名
	return ecdsa.Verify(publicKey, txHash[:], sig.R, sig.S)
}

// 序列化交易
func serialize(tx *Transaction) []byte {

	return []byte(fmt.Sprintf("%s%s%d", tx.Sender, tx.Recipient, tx.Value))
}
