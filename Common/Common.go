package Common

import (
	"github.com/fatih/color"
	"linxBlockchain/Block"
	"linxBlockchain/BlockChain"
	"linxBlockchain/Transaction"
	"linxBlockchain/Wallet"
	"linxBlockchain/util"
	"log"
	"math/big"
	"strconv"
	"time"
)

// DifficultFactor 困难调整因子
var DifficultFactor = 0.2

const MinerAllowance = 100000
const MinerValue = 50000

// Base58Alphabet Base58 编码表
const Base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

var Blockchain *BlockChain.Blockchain
var MINER string
var PasswdWalletMap map[string]*Wallet.Wallet
var KeyWalletMap map[string]*Wallet.Wallet
var TransactionMap map[string]*Transaction.Transaction

func init() {
	wallets, err := util.ReadWallet()
	if err != nil {
		log.Fatal(err)
		return
	}

	bc, err := util.ReadBlockChain()
	if err != nil {
		log.Fatal(err)
		return
	}

	KeyWalletMap = make(map[string]*Wallet.Wallet)
	PasswdWalletMap = make(map[string]*Wallet.Wallet)
	TransactionMap = make(map[string]*Transaction.Transaction)
	creatTime := strconv.FormatInt(time.Now().Unix(), 10)
	if bc != nil {
		Blockchain = &BlockChain.Blockchain{
			Chain:    bc.Chain,
			Balances: bc.Balances,
		}
		for _, wallet := range wallets {
			KeyWalletMap[wallet.PublicKeyHex] = wallet
			PasswdWalletMap[wallet.PassWord] = wallet
		}
		MINER = wallets[0].PublicKeyHex
	} else {
		Blockchain = &BlockChain.Blockchain{
			Chain: []*Block.Block{
				{
					Index:        0,
					Timestamp:    creatTime,
					Transactions: []*Transaction.Transaction{},
					PrevHash:     "000000000000000000000000000000000000",
					Hash:         "000000000000000000000000000000000000",
				},
			},
			Balances: map[string]*big.Int{},
		}
		wallet := NewWallet(Blockchain, strconv.Itoa(int(time.Now().Unix()))[5:])
		Blockchain.Balances[wallet.PublicKeyHex] = big.NewInt(MinerAllowance)
		MINER = wallet.PublicKeyHex
	}
}

// Common========================================================================================================

func ChangeDifficult(avgTime time.Duration) {
	if avgTime <= time.Second*10 {
		DifficultFactor += 0.4
	} else if avgTime < time.Second*20 && avgTime > time.Second*10 {
		DifficultFactor += 0.2
	} else if avgTime > time.Second*20 {
		DifficultFactor *= 0.25
	} else if avgTime > time.Second*40 {
		DifficultFactor *= 0.5
	}

	Block.DIFFICULT = int(float64(Block.DIFFICULT) * (1 + DifficultFactor))
	color.Red("修改难度中>>>\n难度为:%d\n", Block.DIFFICULT)
}
