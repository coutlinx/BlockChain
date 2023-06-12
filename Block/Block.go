package Block

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"linxBlockchain/Transaction"
	"math/big"
)

var DIFFICULT = 255
var IsContinue = true
var ComputTimes = 0

type Block struct {
	Index        int                        `json:"Index"`
	Timestamp    string                     `json:"Timestamp"`
	Transactions []*Transaction.Transaction `json:"Transactions"`
	PrevHash     string                     `json:"PrevHash"`
	Hash         string                     `json:"Hash"`
	Nonce        int                        `json:"Nonce"`
}

func (block *Block) String() string {
	tstr := ""
	for _, ts := range block.Transactions {
		tstr += ts.String()
	}
	return fmt.Sprintf("Index: %d\nTimestamp: %s\nPrevHash: %s\nHash: %s\nNonce:%d\nTransactions:\n%s",
		block.Index, block.Timestamp, block.PrevHash, block.Hash, block.Nonce, tstr)
}

func (block *Block) CalculateHash() {
	record := fmt.Sprintf("%d%s%v%s%d", block.Index, block.Timestamp, block.Transactions, block.PrevHash, block.Nonce)
	hash := sha256.Sum256([]byte(record))
	block.Hash = hex.EncodeToString(hash[:])
}

func (block *Block) AddTransaction(transaction *Transaction.Transaction) {
	block.Transactions = append(block.Transactions, transaction)
	block.IsReadyForMine()
}
func (block *Block) IsReadyForMine() {
	if len(block.Transactions) != 0 {
		block.mineBlock()
	}
}

// MineBlock 挖矿创建新区块
func (block *Block) mineBlock() {
	for {
		if !IsContinue {
			return
		}
		block.CalculateHash()
		ComputTimes++
		if block.vailMine() { // 符合难度目标条件
			//给矿工奖励
			break
		} else {
			block.Nonce++ // 尝试下一个nonce值
		}
	}
}

func (block *Block) vailMine() bool {
	target := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(int64(510)), nil)
	target = target.Mul(target, big.NewInt(15))
	target = target.Div(target, big.NewInt(10))
	if new(big.Int).SetBytes([]byte(block.Hash)).Cmp(target) >= 0 {
		return true
	} else {
		return false
	}
}
