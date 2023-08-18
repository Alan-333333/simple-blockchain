package p2p

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"

	blockchain "github.com/Alan-333333/simple-blockchain/block/chain"
	"github.com/Alan-333333/simple-blockchain/transaction"
)

// 消息类型
const (
	MsgTypeVersion = 1
	MsgTypeTx      = 2
	MsgTypeBlock   = 3
	MsgTypeWallet  = 4
	MsgTypePing    = 999
)

// 版本消息
type Version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

// 编码版本消息
func EncodeVersion(version Version) []byte {
	// 1. 定义buffer
	// 2. 编码Version字段
	// 3. 编码BestHeight字段
	// 4. 编码AddrFrom字符串
	// 5. 返回编码后的内容

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, version.Version)
	binary.Write(buf, binary.LittleEndian, version.BestHeight)

	var addrBytes []byte
	if len(version.AddrFrom) > 0 {
		addrBytes = []byte(version.AddrFrom)
	}

	binary.Write(buf, binary.LittleEndian, uint64(len(addrBytes)))
	buf.Write(addrBytes)

	return buf.Bytes()
}

// 解码版本消息
func DecodeVersion(data []byte) Version {
	// 1. 创建Version对象
	// 2. 解码Version字段
	// 3. 解码BestHeight字段
	// 4. 解码AddrFrom字符串
	// 5. 返回解码后的Version对象
	buf := bytes.NewBuffer(data)

	var version Version

	binary.Read(buf, binary.LittleEndian, &version.Version)
	binary.Read(buf, binary.LittleEndian, &version.BestHeight)

	var addrLen uint64
	binary.Read(buf, binary.LittleEndian, &addrLen)

	addrBytes := make([]byte, addrLen)
	buf.Read(addrBytes)
	version.AddrFrom = string(addrBytes)

	return version
}

// 封装网络消息
type Message struct {
	MsgType int
	Data    []byte
}

// 封装消息编码
func EncodeMessage(msgType int, data interface{}) []byte {
	// 1. 定义buffer
	// 2. 编码消息类型
	// 3. 对data进行编码
	// 4. 返回编码后的消息内容
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, int64(msgType))

	encoder := gob.NewEncoder(buf)
	encoder.Encode(data)

	return buf.Bytes()
}

// 解码消息内容
func DecodeMessage(data []byte) *Message {
	// 1. 创建消息对象
	// 2. 从data中解码消息类型
	// 3. 根据类型解码data字段
	// 4. 返回消息对象
	buf := bytes.NewBuffer(data)

	var msgTypeInt int64
	binary.Read(buf, binary.LittleEndian, &msgTypeInt)

	msg := &Message{MsgType: int(msgTypeInt)}

	decoder := gob.NewDecoder(buf)
	decoder.Decode(&msg.Data)

	return msg
}

// 解码Transaction
func DecodeTransaction(data []byte) (*transaction.Transaction, error) {

	var tx transaction.Transaction
	err := json.Unmarshal(data, &tx)
	return &tx, err
}

// 解码Block
func DecodeBlock(data []byte) (*blockchain.Block, error) {

	var block blockchain.Block
	err := json.Unmarshal(data, &block)

	return &block, err

}
