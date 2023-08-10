package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00) //版本号
const addressChecksumLen = 4

func GenerateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {

	// 1. 生成Curve参数
	c := elliptic.P256()

	// 2. 生成私钥
	privKey, err := ecdsa.GenerateKey(c, rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// 3. 从私钥导出公钥
	pubKey := privKey.PublicKey

	return privKey, &pubKey
}

func GetPublicKey(privateKey *ecdsa.PrivateKey) *ecdsa.PublicKey {

	return &privateKey.PublicKey
}

func Sign(data []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {

	hashed := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashed[:])
	if err != nil {
		return nil, err
	}

	return serializeSig(r, s)
}

func serializeSig(r, s *big.Int) ([]byte, error) {
	var sig bytes.Buffer

	// Write R & S to signature
	err := writeBigInt(&sig, r)
	if err != nil {
		return nil, err
	}

	err = writeBigInt(&sig, s)
	if err != nil {
		return nil, err
	}

	return sig.Bytes(), nil
}

func writeBigInt(w io.Writer, i *big.Int) error {
	// Convert to bytes
	bytes := i.Bytes()

	// Write byte length
	err := binary.Write(w, binary.BigEndian, uint32(len(bytes)))
	if err != nil {
		return err
	}

	// Write bytes
	_, err = w.Write(bytes)
	return err
}

func ParseSig(sigStr string) (r, s *big.Int) {

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

// PubKeyToAddress 将公钥转换为地址
func PubKeyToAddress(pubKey *ecdsa.PublicKey) string {

	// 1. 构造公钥字节数组
	pubBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)

	// 2. SHA256哈希
	hash := sha256.Sum256(pubBytes)

	// 3. 取前20字节作为地址
	address := hash[0:20]

	// 4. 转换为16进制字符串
	return fmt.Sprintf("%x", address)

}

func ParsePrivateKey(keyPEM string) (*ecdsa.PrivateKey, error) {

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

func ParsePubKey(pubKeyPEM string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubKeyPEM))
	if block == nil {
		return nil, errors.New("invalid PEM format")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	pubKey, ok := pubInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not a valid ECDSA public key")
	}

	return pubKey, nil
}

// PubKeyToAddr 将公钥转换为地址
func PubKeyToAddr(pubKey *ecdsa.PublicKey) string {

	// 1. 序列化公钥
	pubKeyBytes := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)

	// 2. 双哈希
	ripmd160 := HashPubKey(pubKeyBytes)

	// 3. 构造版本号和校验和
	payload := append([]byte{version}, ripmd160...)
	checksum := Checksum(payload)

	// 4. 拼接完整数据
	fullPayload := append(payload, checksum...)

	// 5. Base58编码
	addressByte := Base58Encode(fullPayload)

	// 6. addressByte To string
	address := string(addressByte)
	return address

}

// AddrToPubKey 将地址解码为公钥
func AddrToPubKey(addr string) (*ecdsa.PublicKey, error) {

	// 1. Base58解码
	pubKeyHash := Base58Decode([]byte(addr))

	// 2. 分离校验和
	checksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	payload := pubKeyHash[:len(pubKeyHash)-addressChecksumLen]

	// 3. 验证校验和
	if !ValidateChecksum(payload, checksum) {
		return nil, errors.New("invalid checksum")
	}

	// 4. 构造公钥
	pubKey := payload[1:] // 去掉版本号

	if len(pubKey) < 64 {
		return nil, errors.New("invalid pubkey")
	}

	x := pubKey[:32]
	y := pubKey[32:]

	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetBytes(x),
		Y:     new(big.Int).SetBytes(y),
	}, nil

}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

// ValidateChecksum 验证校验和
func ValidateChecksum(payload []byte, checksum []byte) bool {

	// 1. 从payload中提取版本号
	version := payload[0]

	// 2. 构造完整数据
	fullPayload := append([]byte{version}, payload[1:]...)

	// 3. 计算校验和
	expectedChecksum := Checksum(fullPayload)

	// 4. 对比校验和
	return bytes.Equal(expectedChecksum, checksum)

}

func Checksum(payload []byte) []byte {

	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}
