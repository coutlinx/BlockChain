package util

import (
	"encoding/json"
	"github.com/fatih/color"
	"io/ioutil"
)

// SaveMiner 将每个节点的矿工保存在本地,命名为Miner.json
func SaveMiner(port, port2, Miner, Miner2 string) {
	filePath := "D:\\linxBlockchain\\Miner.json"
	var data []byte
	var MinerMap = make(map[string]string)
	MinerMap[port] = Miner
	exists, err := PathExists(filePath)
	if err != nil {
		color.Red(err.Error())
		return
	}
	if exists {
		return
	} else {
		
		if port2 == "" {
			data, err = json.Marshal(MinerMap)
			if err != nil {
				color.Red(err.Error())
				return
			}
		}
	}

	if port2 != "" {
		MinerMap[port2] = Miner2
		data, err = json.Marshal(MinerMap)
		if err != nil {
			color.Red(err.Error())
			return
		}
	}

	SaveFile(filePath, data)
	color.Green("保存矿工到本地中...")
}

// ReadMiner 读取本地Miner.json转换为矿工地址
func ReadMiner() map[string]string {

	var MinerMap = make(map[string]string)
	file, err := ioutil.ReadFile("Miner.json")
	if err != nil {
		color.Green("该节点没有矿工,正在生成...")
		return nil
	}
	err = json.Unmarshal(file, &MinerMap)
	if err != nil {
		color.Red(err.Error())
		return nil
	}
	return MinerMap

}
