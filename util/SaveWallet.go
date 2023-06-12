package util

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"linxBlockchain/Wallet"
)

// SaveWallet 将每个用户创建的钱包保存在本地,命名为Wallet.json
func SaveWallet(w []*Wallet.Wallet) {
	if w == nil {
		return
	}
	data, err := MarshalJSON(w)
	if err != nil {
		fmt.Println(err)
		return
	}
	filePath := "D:\\linxBlockchain\\Wallet.json"
	SaveFile(filePath, data)

	color.Green("保存钱包到本地中...")
}

// ReadWallet 读取本地Wallet.json转换为钱包结构体
func ReadWallet() ([]*Wallet.Wallet, error) {
	color.HiGreen("读取本地钱包文件中...")
	file, err := ioutil.ReadFile("Wallet.json")
	if err != nil {
		color.Green("该节点没有钱包,正在生成...")
		return nil, nil
	}

	var w []*WalletJSON
	var wallets []*Wallet.Wallet
	err = json.Unmarshal(file, &w)
	if err != nil {
		return nil, err
	}
	for _, walletJSON := range w {
		wallet := GetWallet(walletJSON)
		wallets = append(wallets, &wallet)
	}
	return wallets, nil
}
