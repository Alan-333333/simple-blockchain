package p2p

import (
	"testing"

	"github.com/Alan-333333/simple-blockchain/transaction"
	"github.com/Alan-333333/simple-blockchain/wallet"
)

func TestEncodeMessage(t *testing.T) {

	data := "hello"
	msgType := 1

	encoded := EncodeMessage(msgType, []byte(data))

	// 解码编码后的数据

	decodedMsg := DecodeMessage(encoded)

	// 检查解码后的消息类型
	if decodedMsg.MsgType != msgType {
		t.Errorf("Decoded MsgType %d not match input %d", decodedMsg.MsgType, msgType)
	}

	// 检查解码后的消息数据

	if string(decodedMsg.Data) != data {
		t.Errorf("Decoded Data %s not match input %s", decodedMsg.Data, data)
	}

}

func TestDecodeTransaction(t *testing.T) {

	// 1. 创建钱包
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	// 2. 构造交易
	tx := transaction.NewTransaction(walletA.GetAddress(), walletB.GetAddress(), 10)

	// 3. 签名
	tx.Sign(walletA.PrivateKey)
	// 4. 广播交易

	// 编码
	data := EncodeMessage(MsgTypeTx, tx)

	// 调用解码
	decoded, err := DecodeTransaction(data)

	// 检查错误
	if err != nil {
		t.Fatal(err)
	}

	// 检查解码结果
	if decoded.Sender != tx.Sender {
		t.Errorf("Decoded Sender mismatch")
	}

	if decoded.Recipient != tx.Recipient {
		t.Errorf("Decoded Recipient mismatch")
	}

	// 校验其他字段

}
