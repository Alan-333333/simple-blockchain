package wallet

import (
	"testing"
)

func TestEncodedWallet(t *testing.T) {

	// 创建测试钱包
	wallet := NewWallet()

	// 序列化
	data, err := EncodedWallet(wallet)
	if err != nil {
		t.Errorf("encode wallet failed: %v", err)
	}

	// 校验序列化数据长度是否正确
	if len(data) == 0 {
		t.Error("encoded wallet data is empty")
	}

	// 反序列化
	restored, err := DecodeWallet(data)
	if err != nil {
		t.Errorf("decode wallet failed: %v", err)
	}

	// 校验反序列化后钱包属性
	if restored.Address != wallet.Address {
		t.Errorf("address not matched, got %s want %s", restored.Address, wallet.Address)
	}

}
