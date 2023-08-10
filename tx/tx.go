package tx

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
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
func (tx *Transaction) Sign(privateKeyStr string) error {
	// 将私钥字符串解析为私钥对象
	privateKey, err := parsePrivateKey(privateKeyStr)
	if err != nil {
		return err
	}
	// 签名算法
	signature := signECDSA(privateKey, tx) // 使用私钥签名

	tx.Signature = signature // 添加签名到交易

	return nil
}

// isValid 验证交易签名
func (tx *Transaction) isValid() bool {
	// 解析公钥
	pubKey, err := parsePubKey(tx.Sender)
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
	r, s := parseSig(sig)

	// 序列化交易
	txBytes := serialize(tx)

	// 计算交易哈希
	txHash := sha256.Sum256(txBytes)

	// 使用公钥验证签名
	return ecdsa.Verify(publicKey, txHash[:], r, s)

}

func parsePrivateKey(keyPEM string) (*ecdsa.PrivateKey, error) {

	// 从pem格式块解析出DER编码的私钥
	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	// 解析DER编码的私钥
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func parsePubKey(pubKeyStr string) (*ecdsa.PublicKey, error) {

	// 从pem格式块解析出DER编码的私钥
	block, _ := pem.Decode([]byte(pubKeyStr))
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	// 从字符串反序列化公钥
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pubKey.(*ecdsa.PublicKey), nil
}

func serialize(tx *Transaction) []byte {
	// 序列化交易为字节数组
	return []byte(fmt.Sprintf("%v%v", tx))
}

func parseSig(sigStr string) (r, s *big.Int) {

	// 将签名的16进制字符串转为字节数组
	sig, err := hex.DecodeString(sigStr)
	if err != nil {
		return
	}

	// 签名由r和s组成,均32字节
	sigLen := len(sig)
	if sigLen != 64 {
		return
	}

	// 提取r和s
	rBytes := sig[:32]
	sBytes := sig[32:]

	// 转换为big.Int类型
	r = new(big.Int).SetBytes(rBytes)
	s = new(big.Int).SetBytes(sBytes)

	return
}
