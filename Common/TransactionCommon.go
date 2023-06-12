package Common

import (
	"crypto/ecdsa"
	"errors"
	"github.com/fatih/color"
	"linxBlockchain/Sign"
	"linxBlockchain/Transaction"
	"linxBlockchain/Wallet"
	"math/big"
)

func CreateTransaction(sender, receiver string, amount *big.Int) (*Transaction.Transaction, error) {
	senderWallet := GetWalletByPublicKey(sender)
	tx := &Transaction.Transaction{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}

	SingForTransaction(senderWallet, tx)

	flag, err := isVail(sender, tx)
	if err != nil {
		return nil, err
	}
	if flag {
		tx.CalculateHash()
		return tx, nil
	}
	return nil, errors.New("未知错误创建交易或签名")

}

func sending(sender, receiver string, amount *big.Int) {
	color.Yellow("%d\n", Blockchain.Balances[receiver])
	if sender != "System" {
		Blockchain.Balances[sender].Sub(Blockchain.Balances[sender], amount)
		Blockchain.Balances[receiver].Add(Blockchain.Balances[receiver], amount)
	} else {
		Blockchain.Balances[receiver].Add(Blockchain.Balances[receiver], amount)
	}

}

func isVail(publickey string, tx *Transaction.Transaction) (bool, error) {
	value := Blockchain.Balances[publickey]
	wallet := GetWalletByPublicKey(publickey)
	if value == nil {
		return false, errors.New("不存在发送方")
	}
	result := value.Cmp(tx.Amount)
	if result <= 0 {
		return false, errors.New("发送方没有那么多钱")
	}
	if !VailTransaction(wallet, tx.Sign) {
		return false, errors.New("签名错误")
	}
	return true, nil
}

func VailTransaction(wallet *Wallet.Wallet, sign *Sign.Sign) bool {
	return ecdsa.Verify(wallet.GetPublicKey(), sign.Message, sign.R, sign.S)
}

func GetAllTransactions() []*Transaction.Transaction {
	var trans []*Transaction.Transaction
	for _, transaction := range TransactionMap {
		trans = append(trans, transaction)
	}
	return trans
}
