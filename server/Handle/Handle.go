package Handle

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"linxBlockchain/Block"
	"linxBlockchain/Common"
	"linxBlockchain/Transaction"
	"linxBlockchain/util"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

// Transactions 新增交易,创建交易并且分享给好兄弟节点
func Transactions(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	switch request.Method {
	case http.MethodPost:
		amount, err := strconv.ParseInt(request.FormValue("Amount"), 10, 64)
		var t = &Transaction.Transaction{
			Sender:   request.FormValue("Sender"),
			Receiver: request.FormValue("Receiver"),
			Amount:   big.NewInt(amount),
		}

		if err != nil {
			color.Red(err.Error())
			return
		}

		t.Amount = big.NewInt(amount)
		if err != nil {
			color.Red(err.Error())
			return
		}
		if t == nil {
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				color.Red(err.Error())
				return
			}

			err = json.Unmarshal(body, &t)
			if err != nil {
				color.Red(err.Error())
				return
			}

		}
		t.CalculateHash()
		// 序列化transaction用作传参
		jsonData, err := json.Marshal(t)
		if err != nil {
			fmt.Println("Failed to marshal JSON data:", err)
			return
		}
		go transactionHashReturn(writer, t)
		go sendBrother(jsonData, *t)
		if Block.IsContinue {
			color.Blue(strconv.Itoa(Port) + "开始挖矿咯!")
			// 通知自己该开始挖矿啦！
			startMine(*t)
		}

	}

}

// GetTransactions 接收交易并然后开始挖矿
func GetTransactions(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	switch request.Method {
	case http.MethodPost:
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			color.Red(err.Error())
			return
		}
		var t Transaction.Transaction
		err = json.Unmarshal(body, &t)
		if err != nil {
			color.Red(err.Error())
			return
		}
		color.Cyan("获得好兄弟给的交易\n%s", t.String())
		if Block.IsContinue {
			color.Blue(strconv.Itoa(Port) + "开始挖矿咯!")
			// 通知自己该开始挖矿啦！
			//time.Sleep(time.Second)
			startMine(t)
		}

	}
}

// FinishMining 结束挖矿,修改Block包里面的变量从而停止挖矿.
func FinishMining(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	switch request.Method {
	case http.MethodGet:
		color.HiCyan("挖矿慢了呀,兄弟节点挖矿成功!")
		Block.IsContinue = false
		syncBlock(brotherNode)
	}

}

// FindYou 找到好兄弟然后给他发我自己的端口号
func FindYou(writer http.ResponseWriter, req *http.Request) {
	enableCors(&writer)
	switch req.Method {
	case http.MethodGet:
		IsFind = true
		brotherNode = req.URL.Query().Get("brother")
		color.Cyan("发现好兄弟%s\n", brotherNode)
	}
}

// GetChain 获取本地链
func GetChain(writer http.ResponseWriter, req *http.Request) {
	enableCors(&writer)
	switch req.Method {
	case http.MethodGet:
		writer.Header().Add("Content-Type", "application/json")
		bc := Common.Blockchain
		jsonData, err := json.Marshal(bc)
		if err != nil {
			fmt.Println("JSON serialization error:", err)
			return
		}
		io.WriteString(writer, string(jsonData))
	}
}

// GetBrother 返回兄弟节点的网址
func GetBrother(writer http.ResponseWriter, req *http.Request) {
	enableCors(&writer)
	switch req.Method {
	case http.MethodGet:
		writer.Header().Add("Content-Type", "application/json")
		io.WriteString(writer, string(brotherNode))
	}
}

// GetMiner 获取该节点的矿工公钥
func GetMiner(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	switch request.Method {
	case http.MethodGet:
		io.WriteString(writer, Common.MINER)
	}
}

// GetBalance 获取全部用户的余额
func GetBalance(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	switch request.Method {
	case http.MethodGet:
		Account := request.FormValue("Account")
		io.WriteString(writer, Common.Blockchain.Balances[Account].String())
	}
}

// GetTran 通过交易Hash获取交易
func GetTran(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	switch request.Method {
	case http.MethodGet:
		writer.Header().Add("Content-Type", "application/json")
		Hash := request.FormValue("Hash")
		data, err := json.Marshal(Common.TransactionMap[Hash])
		if err != nil {
			color.Red(err.Error())
		}
		io.WriteString(writer, string(data))
	}
}

// GetBlock 通过Hash或者区块获取区块全部信息
func GetBlock(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	writer.Header().Add("Content-Type", "application/json")
	switch request.Method {
	case http.MethodGet:

		index := request.FormValue("Index")
		Hash := request.FormValue("Hash")
		//color.Red("Index:%s\n", index)
		//color.Red("Hash:%s\n", Hash)
		if index == "" {
			data, err := json.Marshal(Common.Blockchain.GetBlockByHash(Hash))
			if err != nil {
				color.Red(err.Error())
			}
			io.WriteString(
				writer,
				string(data),
			)
		} else if Hash == "" {
			BlockIndex, err := strconv.Atoi(index)
			if err != nil {
				color.Red(err.Error())
				return
			}
			if BlockIndex >= len(Common.Blockchain.Chain) {
				io.WriteString(writer, "区块不存在")
			} else {
				data, err := json.Marshal(Common.Blockchain.Chain[BlockIndex])
				if err != nil {
					color.Red(err.Error())
				}

				io.WriteString(
					writer,
					string(data),
				)
			}

		}
	}
}

// AsyncBlock 同步区块链
func AsyncBlock(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	switch request.Method {
	case http.MethodGet:
		syncBlock(brotherNode)
		syncWallet(brotherNode)
	}
}

// OpenMining 开启挖矿
func OpenMining(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	switch request.Method {
	case http.MethodGet:
		color.HiRed("IsContinue:%v\n", Block.IsContinue)
		color.HiMagenta("收到可以开始准备挖矿\n")
		Block.IsContinue = true
	}
}

// GetWallet 全部用户的钱包
func GetWallet(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	switch request.Method {
	case http.MethodGet:
		writer.Header().Add("Content-Type", "application/json")
		wm := WalletMap{
			KeyWalletMap:    Common.KeyWalletMap,
			PasswdWalletMap: Common.PasswdWalletMap,
		}
		jsonData, err := json.Marshal(wm)
		if err != nil {
			fmt.Println("JSON serialization error:", err)
			return
		}
		_, err = io.WriteString(writer, string(jsonData))
		if err != nil {
			color.Red(err.Error())
			return
		}
	}
}

// CreatWallet 创建钱包
func CreatWallet(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	writer.Header().Add("Content-Type", "application/json")
	switch request.Method {
	case http.MethodPost:
		Password := request.PostFormValue("Password")
		w := newWallet(Password)
		data, err := json.Marshal(w)
		if err != nil {
			color.Red(err.Error())
		}
		_, err = http.Get(brotherNode + "/asyncBlock")
		if err != nil {
			color.Red(err.Error())
			return
		}
		color.HiCyan("用户创建钱包中...")
		io.WriteString(writer, string(data))
	}
}

// LoadWallet 加载钱包
func LoadWallet(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	writer.Header().Add("Content-Type", "application/json")
	switch request.Method {
	case http.MethodPost:
		file, _, err := request.FormFile("file")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		wallet := util.GetWalletFromFile(file)
		addWallet(wallet)
		_, err = http.Get(brotherNode + "/asyncBlock")
		if err != nil {
			color.Red(err.Error())
			return
		}
		data, err := json.Marshal(wallet)
		if err != nil {
			color.Red(err.Error())
		}

		io.WriteString(writer, string(data))
	}

}

func GetAllTransactions(writer http.ResponseWriter, request *http.Request) {
	enableCors(&writer)
	switch request.Method {
	case http.MethodGet:
		writer.Header().Add("Content-Type", "application/json")
		trans := getAllTransaction()
		//color.Red("%v\n", trans)
		data, err := json.Marshal(trans)
		if err != nil {
			return
		}
		io.WriteString(writer, string(data))
	}
}
