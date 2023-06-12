package Transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"linxBlockchain/Sign"
	"math/big"
	"time"
)

type Transaction struct {
	Sender   string     `json:"Sender,omitempty"`
	Receiver string     `json:"Receiver,omitempty"`
	Amount   *big.Int   `json:"Amount,omitempty"`
	Hash     string     `json:"Hash,omitempty"`
	Sign     *Sign.Sign `json:"Sign,omitempty"`
}

// CalculateHash 计算交易的哈希值
func (tx *Transaction) CalculateHash() {
	record := fmt.Sprintf("%s%s%d%d%v", tx.Sender, tx.Receiver, tx.Amount, time.Now().Unix(), tx.Sign)
	hash := sha256.Sum256([]byte(record))
	tx.Hash = hex.EncodeToString(hash[:])
}

func (tx *Transaction) String() string {
	return fmt.Sprintf("Sender:%s\tReceiver:%s\tAmount:%d\tHash:%s\t", tx.Sender, tx.Receiver, tx.Amount, tx.Hash)
}
