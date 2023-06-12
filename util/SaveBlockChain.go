package util

import (
	"encoding/json"
	"github.com/fatih/color"
	"io/ioutil"
	"linxBlockchain/BlockChain"
)

// Save 将区块链保存在本地,命名为BlockChain.json
func Save(bc *BlockChain.Blockchain) error {
	data, err := json.Marshal(bc)
	if err != nil {
		color.Red(err.Error())
		return err
	}
	filePath := "D:\\linxBlockchain\\BlockChain.json"
	SaveFile(filePath, data)
	color.Green("保存区块链到本地中...")
	return nil
}

// ReadBlockChain 读取本地BlockChain.json转换为区块链结构体
func ReadBlockChain() (*BlockChain.Blockchain, error) {
	color.HiGreen("读取本地区块链文件中...")

	file, err := ioutil.ReadFile("BlockChain.json")
	if err != nil {
		color.Green("该节点没有创世纪块,正在生成...")
		return nil, nil
	}

	var bc *BlockChain.Blockchain
	err = json.Unmarshal(file, &bc)
	if err != nil {
		return nil, err
	}
	return bc, nil

}
