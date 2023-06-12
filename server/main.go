package main

import (
	"flag"
	"fmt"
	"linxBlockchain/server/Handle"
	"log"
	"net/http"
	"strconv"
)

func main() {
	flag.IntVar(&Handle.Port, "port", 5000, "端口号默认为5000")
	flag.Parse()
	fmt.Printf("port::%d \n", Handle.Port)
	//CreatServer()
	go Handle.CloseHandle()
	go Handle.StartHandle()
	http.HandleFunc("/", Handle.GetChain)
	http.HandleFunc("/brother", Handle.GetBrother)
	http.HandleFunc("/transactions", Handle.Transactions)
	http.HandleFunc("/getTransactions", Handle.GetTransactions)
	http.HandleFunc("/find", Handle.FindYou)
	http.HandleFunc("/FinishMining", Handle.FinishMining)
	http.HandleFunc("/wallet", Handle.GetWallet)
	http.HandleFunc("/openMining", Handle.OpenMining)
	http.HandleFunc("/getBlock", Handle.GetBlock)
	http.HandleFunc("/getTran", Handle.GetTran)
	http.HandleFunc("/getBalance", Handle.GetBalance)
	http.HandleFunc("/getMiner", Handle.GetMiner)
	http.HandleFunc("/asyncBlock", Handle.AsyncBlock)
	http.HandleFunc("/creatWallet", Handle.CreatWallet)
	http.HandleFunc("/loadWallet", Handle.LoadWallet)
	http.HandleFunc("/getAllTransaction", Handle.GetAllTransactions)
	fmt.Println(strconv.Itoa(Handle.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(Handle.Port), nil))
}
