package Wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

type Wallet struct {
	PublicKeyHex  string
	PrivateKeyHex string
	BlockAddress  string
	publicKey     *ecdsa.PublicKey
	privateKey    *ecdsa.PrivateKey
	PassWord      string
}

// GenerateKeyPairFromPassword 生成公钥和私钥
func GenerateKeyPairFromPassword() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {

	curve := elliptic.P256()

	// 生成私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// 获取公钥
	publicKey := &privateKey.PublicKey

	return privateKey, publicKey, nil
}
func (wallet *Wallet) SetPublicKey(publicKey *ecdsa.PublicKey) {
	wallet.publicKey = publicKey
}

func (wallet *Wallet) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	wallet.privateKey = privateKey
}

func (wallet *Wallet) GetPublicKey() *ecdsa.PublicKey {
	return wallet.publicKey
}

func (wallet *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	return wallet.privateKey
}

func (wallet *Wallet) String() string {
	return fmt.Sprintf("PublicKeyHex: [%s]\n PrivateKeyHex: [%s]\n BlockAddress: [%s]", wallet.PublicKeyHex, wallet.PrivateKeyHex, wallet.BlockAddress)
}
