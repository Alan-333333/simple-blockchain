package p2p

import (
	"testing"
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
