package Handle

import (
	"linxBlockchain/Wallet"
	"time"
)

// Port 本身的端口号
var Port int

// 兄弟节点的路径
var brotherNode string

// IsFind 是否找到兄弟节点
var IsFind = false

// StartTime 程序运行的时间(用于修改难度)
var StartTime time.Time

// PortRange 寻找节点的范围
const PortRange = 200

// WalletMap 同步节点之间的钱包用的中间结构体
type WalletMap struct {
	PasswdWalletMap map[string]*Wallet.Wallet
	KeyWalletMap    map[string]*Wallet.Wallet
}
type JsonWallet struct {
	PrivateKey string `json:"PrivateKey"`
	Password   string `json:"Password"`
}
