package Common

import (
	"linxBlockchain/Block"
	"linxBlockchain/Transaction"
	"math/big"
	"strconv"
	"time"
)

func NewBlock(index int, transaction []*Transaction.Transaction, prevHash string, nonce ...int) *Block.Block {
	creatTime := strconv.FormatInt(time.Now().Unix(), 10)
	block := &Block.Block{
		Index:        index,
		Timestamp:    creatTime,
		Transactions: transaction,
		PrevHash:     prevHash,
		Nonce:        nonce[0],
	}
	block.IsReadyForMine()
	return block
}

func AddBlock(Miner string, block *Block.Block) {
	Blockchain.AddBlock(block)

	for _, tx := range block.Transactions {
		TransactionMap[tx.Hash] = tx
		sending(tx.Sender, tx.Receiver, tx.Amount)
	}
	sending("System", Miner, big.NewInt(MinerValue))
}
