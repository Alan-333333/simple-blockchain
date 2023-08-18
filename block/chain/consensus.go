package blockchain

type Consensus interface {
	GenerateBlock(block *Block)
	VerifyBlock(block *Block) bool
}
