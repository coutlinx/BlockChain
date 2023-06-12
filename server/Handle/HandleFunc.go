package Handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"linxBlockchain/Block"
	"linxBlockchain/BlockChain"
	"linxBlockchain/Common"
	"linxBlockchain/Transaction"
	"linxBlockchain/Wallet"
	"linxBlockchain/util"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"time"
)

// SyncBlockForCreat 在程序启动的时候跟其他节点进行同步
func SyncBlockForCreat(nodeAddress string) {
	response, err := http.Get(nodeAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var bc BlockChain.Blockchain
	err = json.Unmarshal(body, &bc)

	// 打印响应内容
	timestamp1, err := strconv.ParseInt(bc.Chain[0].Timestamp, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	timestamp2, err := strconv.ParseInt(Common.Blockchain.Chain[0].Timestamp, 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	color.Yellow("同步区块中...")

	if timestamp1 <= timestamp2 {
		Common.Blockchain.Chain = bc.Chain
	}

	if !reflect.DeepEqual(Common.Blockchain.Balances, bc.Balances) {
		for key, value := range bc.Balances {
			Common.Blockchain.Balances[key] = value
		}
	}
}

// ChangeDifficult 修改pow的难度
func ChangeDifficult() {
	for {
		time.AfterFunc(time.Second*60, func() {
			costTime := time.Now().Sub(StartTime)
			nanoseconds := costTime.Nanoseconds()
			avgTime := time.Duration(nanoseconds/int64(len(Common.Blockchain.Chain))) * time.Nanosecond
			Common.ChangeDifficult(avgTime)
		})
	}
}

// StartHandle 启动程序的时候的操作(绑定矿工)
func StartHandle() {
	StartTime = time.Now()
	if Common.MINER == "" {
		Miner := util.ReadMiner()["http://localhost"+strconv.Itoa(Port)]
		if Miner != "" {
			Common.MINER = Miner
		}
	}

	for _, block := range Common.Blockchain.Chain {
		if len(block.Transactions) == 0 {
			continue
		}
		Common.TransactionMap[block.Transactions[0].Hash] = block.Transactions[0]
	}

	searchNodes()
	fmt.Println("======================矿工信息====================")
	fmt.Println(Common.GetWalletByPublicKey(Common.MINER))
	fmt.Println("======================创世纪块====================")
	fmt.Println(Common.Blockchain.Chain[0])
}

// CloseHandle 关闭程序的时候做的操作(区块链本地化)
func CloseHandle() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	// 收到关闭信号后保存区块链数据
	var wallets []*Wallet.Wallet
	for _, w := range Common.KeyWalletMap {
		wallets = append(wallets, w)
	}
	util.SaveWallet(wallets)
	err := util.Save(Common.Blockchain)
	if err != nil {
		fmt.Println("保存区块链失败:", err)
	}
	SaveMiner()
	color.Magenta("关闭节点中...")
	os.Exit(0)
}

// 开始挖矿,创建区块，把区块添加到链上
func startMine(t Transaction.Transaction) {

	// 添加区块到区块链成功后就通知兄弟节点说别挖啦,哥们已经挖完了.
	block := creatBlock([]*Transaction.Transaction{&t})

	if block.Hash != "" {
		Common.AddBlock(Common.MINER, block)
		color.Green(strconv.Itoa(Port) + "挖矿成功！")
		_, err := http.Get(brotherNode + "/FinishMining")
		if err != nil {
			color.Red(err.Error())
			return
		}
	}

}

// 跨域解决
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// 给兄弟节点发送交易
func sendBrother(jsonData []byte, t Transaction.Transaction) {
	color.Cyan("发给兄弟节点说有交易啦!\n%s", t.String())
	// 发给兄弟节点说有交易啦!
	_, err := http.Post(brotherNode+"/getTransactions", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		color.Red(err.Error())
		return
	}
}

// 给调用交易的用户返回Hash
func transactionHashReturn(w http.ResponseWriter, tx *Transaction.Transaction) {
	io.WriteString(w, tx.String()+"Hash:"+tx.Hash)
}

// 同步节点的区块
func syncBlock(nodeAddress string) {
	response, err := http.Get(nodeAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var bc BlockChain.Blockchain
	err = json.Unmarshal(body, &bc)

	brotherTimestamp, err := strconv.ParseInt(bc.Chain[len(bc.Chain)-1].Timestamp, 10, 64)
	if err != nil {
		color.Red(err.Error())
		return
	}

	selfTimeStamp, err := strconv.ParseInt(Common.Blockchain.Chain[len(Common.Blockchain.Chain)-1].Timestamp, 10, 64)
	if err != nil {
		color.Red(err.Error())
		return
	}
	color.Blue("同步兄弟节点中...")

	if len(Common.Blockchain.Chain) < len(bc.Chain) {
		//fmt.Println(123)
		Common.Blockchain.Chain = append(Common.Blockchain.Chain, bc.Chain[len(bc.Chain)-1])
	}

	if brotherTimestamp == selfTimeStamp && len(Common.Blockchain.Chain) == len(bc.Chain) && bc.Chain[len(bc.Chain)-1].Nonce > Common.Blockchain.Chain[len(Common.Blockchain.Chain)-1].Nonce {
		//fmt.Println(1234)
		Common.Blockchain.Chain[len(Common.Blockchain.Chain)-1] = bc.Chain[len(bc.Chain)-1]
		//Common.Blockchain.Chain = bc.Chain
	}

	if brotherTimestamp < selfTimeStamp {
		//fmt.Println(12345)
		Common.Blockchain.Chain[len(Common.Blockchain.Chain)-1] = bc.Chain[len(bc.Chain)-1]
	}

	if !reflect.DeepEqual(Common.Blockchain.Balances, bc.Balances) {
		for key, value := range bc.Balances {
			Common.Blockchain.Balances[key] = value
		}
	}

}

// 同步节点间的钱包
func syncWallet(nodeAddress string) {
	color.HiYellow("同步钱包中...")
	response, err := http.Get(nodeAddress + "/wallet")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var wm WalletMap
	err = json.Unmarshal(body, &wm)
	if err != nil {
		color.Red(err.Error())
		return
	}
	//color.Yellow("KeyWalletMap\n")
	for key, value := range wm.KeyWalletMap {
		//color.Blue("%s:%v\n", key, value)
		Common.KeyWalletMap[key] = value
	}

	//color.Yellow("PasswdWalletMap\n")
	for key, value := range wm.PasswdWalletMap {
		//color.Blue("%s:%v\n", key, value)
		Common.PasswdWalletMap[key] = value
	}

}

// 查找附近的节点
func searchNodes() {
	for {
		if IsFind {
			break
		}
		for i := 0; i < PortRange; i++ {
			if IsFind {
				break
			}

			if 5000+i == Port {
				continue
			}

			nodeAddress := fmt.Sprintf("localhost:%d", 5000+i)
			conn, err := net.DialTimeout("tcp", nodeAddress, time.Second)
			if err == nil {
				brotherNode = "http://" + nodeAddress
				color.Cyan("发现好兄弟%s\n", brotherNode)
				_, err := http.Get(brotherNode + "/find?brother=http://localhost:" + strconv.Itoa(Port))
				if err != nil {
					log.Fatal(err)
				}
				conn.Close()
				IsFind = true
				break
			}

		}
	}
	SyncBlockForCreat(brotherNode)
	syncWallet(brotherNode)
}

// SaveMiner 保存矿工信息
func SaveMiner() {
	client := &http.Client{
		Timeout: 1 * time.Second, // Set timeout duration
	}
	request, err := http.NewRequest("GET", brotherNode+"/getMiner", nil)
	if err != nil {
		color.Magenta("未获取到兄弟节点的矿工,保存中...")
		util.SaveMiner("http://localhost:"+strconv.Itoa(Port), "", Common.MINER, "")
		return
	}

	response, err := client.Do(request)
	if err != nil {
		color.Magenta("未获取到兄弟节点的矿工,保存中...")
		util.SaveMiner("http://localhost:"+strconv.Itoa(Port), "", Common.MINER, "")
		return
	}

	defer response.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	util.SaveMiner("http://localhost:"+strconv.Itoa(Port), brotherNode, Common.MINER, string(body))
}

// 创建区块
func creatBlock(tp []*Transaction.Transaction) *Block.Block {
	rand.Seed(time.Now().UnixNano())
	// 生成随机整数
	randomInt := rand.Intn(1000000000) // 生成0到1亿之间的随机整数

	return Common.NewBlock(len(Common.Blockchain.Chain), tp, Common.Blockchain.Chain[len(Common.Blockchain.Chain)-1].Hash, randomInt)

}

// 创建钱包
func newWallet(Password string) *Wallet.Wallet {
	wallet := Common.NewWallet(Common.Blockchain, Password)
	return wallet
}

// 导入钱包
func addWallet(wallet *Wallet.Wallet) {
	Common.AddWallet(wallet)
	color.HiCyan("导入钱包成功...")
}

// 获得全部的交易信息
func getAllTransaction() []*Transaction.Transaction {
	return Common.GetAllTransactions()
}
