package blockchain

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const chainMetaFile = "./dat/blockchain/meta.json"

type BlockchainMetadata struct {
	LastBlockHash []byte
	BlockCount    int
}

func (m *BlockchainMetadata) Save() error {

	// 1. 序列化为json格式
	metaJson, err := json.Marshal(m)
	if err != nil {
		return err
	}
	// 2. 创建保存目录
	err = os.MkdirAll(filepath.Dir(chainMetaFile), 0700)
	if err != nil {
		return err
	}

	// 2. 写入文件
	return os.WriteFile(chainMetaFile, metaJson, 0600)
}
