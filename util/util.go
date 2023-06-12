package util

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"linxBlockchain/Wallet"
	"math/big"
	"mime/multipart"
	"os"
)

// WalletJSON 方便中间读取钱包的结构体
type WalletJSON struct {
	PublicKeyHex  string `json:"PublicKeyHex"`
	PrivateKeyHex string `json:"PrivateKeyHex"`
	BlockAddress  string `json:"BlockAddress"`
	PassWord      string `json:"PassWord"`
}

// SaveFile 保存文件的公共方法
func SaveFile(fileName string, data []byte) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0766)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	_, err = write.Write(data)
	if err != nil {
		color.Red(err.Error())
		return
	}
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

// MarshalJSON 将钱包的结构体进行转换让它能正常存储进我们的json文件中
func MarshalJSON(w []*Wallet.Wallet) ([]byte, error) {

	var wallets []*WalletJSON

	for _, wallet := range w {
		walletJson := &WalletJSON{
			PublicKeyHex:  wallet.PublicKeyHex,
			PrivateKeyHex: wallet.PrivateKeyHex,
			BlockAddress:  wallet.BlockAddress,
			PassWord:      wallet.PassWord,
		}
		wallets = append(wallets, walletJson)
	}
	return json.Marshal(wallets)
}

// GetKey 将读取出来的privateKey来转换成ecdsa类型的privateKey和publicKey
func getKey(PrivateKeyHex string) (*ecdsa.PublicKey, *ecdsa.PrivateKey) {
	privateKeyD := new(big.Int)
	privateKeyD.SetString(PrivateKeyHex, 16)
	curve := elliptic.P256()
	x, y := curve.ScalarBaseMult(privateKeyD.Bytes())
	publicKey := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
	tempKey := new(ecdsa.PrivateKey)
	tempKey.D = privateKeyD
	return publicKey, tempKey
}
func GetWallet(walletJSON *WalletJSON) Wallet.Wallet {
	publickey, privatekey := getKey(walletJSON.PrivateKeyHex)
	wallet := Wallet.Wallet{
		PublicKeyHex:  walletJSON.PublicKeyHex,
		PrivateKeyHex: walletJSON.PrivateKeyHex,
		BlockAddress:  walletJSON.BlockAddress,
		PassWord:      walletJSON.PassWord,
	}
	wallet.SetPublicKey(publickey)
	wallet.SetPrivateKey(privatekey)
	return wallet
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetWalletFromFile(file multipart.File) *Wallet.Wallet {
	result, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}
	var wj *WalletJSON

	err = json.Unmarshal(result, &wj)
	if err != nil {
		return nil
	}
	wallet := GetWallet(wj)
	return &wallet
}
