package BlockChain

import (
	"fmt"
	"linxBlockchain/Block"
	"math/big"
)

type Blockchain struct {
	Chain    []*Block.Block      `json:"Chain"`
	Balances map[string]*big.Int `json:"Balances"`
}

func (bc *Blockchain) String() string {
	Blocks := ""
	for i, block := range bc.Chain {
		Blocks += fmt.Sprintf("BLOCK:%d", i) + "\n" + block.String() + "\n"
	}
	return fmt.Sprintf("BLOCKCHAIN:\n%s", Blocks)
}
func (bc *Blockchain) GetBlockByHash(Hash string) *Block.Block {
	for _, block := range bc.Chain {
		if block.Hash == Hash {
			return block
		}
	}
	return nil
}

func (bc *Blockchain) AddBlock(block *Block.Block) {
	bc.Chain = append(bc.Chain, block)
}
