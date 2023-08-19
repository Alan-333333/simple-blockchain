package blockchain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVerifyBlock(t *testing.T) {
	// 构造一个有效的区块
	validBlock := CreateGenesisBlock()
	pow := POW{}

	t.Run("valid block", func(t *testing.T) {
		result := pow.VerifyBlock(validBlock)
		assert.True(t, result)
	})

	t.Run("invalid version", func(t *testing.T) {
		invalidBlock := validBlock
		invalidBlock.Version++
		result := pow.VerifyBlock(invalidBlock)
		assert.False(t, result)
	})

	t.Run("invalid timestamp", func(t *testing.T) {
		invalidBlock := validBlock
		invalidBlock.Timestamp = uint64(time.Now().Unix()) + 1000
		result := pow.VerifyBlock(invalidBlock)
		assert.False(t, result)
	})

	t.Run("invalid hash", func(t *testing.T) {
		invalidBlock := validBlock
		invalidBlock.Hash = []byte("invalid hash")
		result := pow.VerifyBlock(invalidBlock)
		assert.False(t, result)
	})

}
