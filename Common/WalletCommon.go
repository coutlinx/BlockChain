package Common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/fatih/color"
	"golang.org/x/crypto/ripemd160"
	"linxBlockchain/BlockChain"
	"linxBlockchain/Sign"
	"linxBlockchain/Transaction"
	"linxBlockchain/Wallet"
	"log"
	"math/big"
)

func NewWallet(blockChain *BlockChain.Blockchain, password string) *Wallet.Wallet {
	PrivateKey, PublicKey, err := Wallet.GenerateKeyPairFromPassword()
	if err != nil {
		log.Fatal(err)
	}
	wallet := &Wallet.Wallet{
		PublicKeyHex:  hex.EncodeToString(PublicKey.X.Bytes()),
		PrivateKeyHex: hex.EncodeToString(PrivateKey.D.Bytes()),
		PassWord:      password,
	}
	wallet.SetPublicKey(PublicKey)
	wallet.SetPrivateKey(PrivateKey)
	address, err := GenerateBlockchainAddress(PublicKey)
	if err != nil {
		color.Red(err.Error())
	}
	wallet.BlockAddress = address

	blockChain.Balances[wallet.PublicKeyHex] = big.NewInt(0)
	color.Green("Balance:%d\n", blockChain.Balances[wallet.PublicKeyHex])
	KeyWalletMap[wallet.PublicKeyHex] = wallet
	PasswdWalletMap[password] = wallet
	return wallet
}
func NewWalletByKey(blockChain *BlockChain.Blockchain, Password string, PublicKey *ecdsa.PublicKey, PrivateKey *ecdsa.PrivateKey) {
	wallet := &Wallet.Wallet{
		PublicKeyHex:  hex.EncodeToString(PublicKey.X.Bytes()),
		PrivateKeyHex: hex.EncodeToString(PrivateKey.D.Bytes()),
		PassWord:      Password,
	}
	wallet.SetPublicKey(PublicKey)
	wallet.SetPrivateKey(PrivateKey)
	address, err := GenerateBlockchainAddress(PublicKey)
	if err != nil {
		color.Red(err.Error())
	}
	wallet.BlockAddress = address
	blockChain.Balances[wallet.PublicKeyHex] = big.NewInt(0)
	KeyWalletMap[wallet.PublicKeyHex] = wallet
	PasswdWalletMap[Password] = wallet
}
func AddWallet(wallet *Wallet.Wallet) {
	Blockchain.Balances[wallet.PublicKeyHex] = big.NewInt(0)
	KeyWalletMap[wallet.PublicKeyHex] = wallet
	PasswdWalletMap[wallet.PassWord] = wallet
}
func SingForTransaction(wallet *Wallet.Wallet, ts *Transaction.Transaction) {
	message := sha256.Sum256([]byte(ts.String()))
	r, s, err := ecdsa.Sign(rand.Reader, wallet.GetPrivateKey(), message[:])
	if err != nil {
		return
	}
	ts.Sign = &Sign.Sign{
		Message: message[:],
		R:       r,
		S:       s,
	}
}
func GetWalletByPasswd(passwd string) *Wallet.Wallet {
	return PasswdWalletMap[passwd]
}

func GetWalletByPublicKey(PublicKey string) *Wallet.Wallet {
	return KeyWalletMap[PublicKey]
}

func GenerateBlockchainAddress(publicKey *ecdsa.PublicKey) (string, error) {
	// 1. 获取公钥的字节形式
	publicKeyBytes := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)

	// 2. 对公钥进行 SHA-256 哈希运算
	sha256Hash := sha256.Sum256(publicKeyBytes)

	// 3. 对哈希结果进行 RIPEMD-160 哈希运算
	ripemd160Hash := ripemd160.New()
	_, err := ripemd160Hash.Write(sha256Hash[:])
	if err != nil {
		return "", err
	}
	ripemd160HashResult := ripemd160Hash.Sum(nil)

	// 4. 添加地址版本前缀
	versionPrefix := []byte{0x00}
	versionPrefixedHash := append(versionPrefix, ripemd160HashResult...)

	// 5. 进行两次 SHA-256 哈希运算
	sha256Hash1 := sha256.Sum256(versionPrefixedHash)
	sha256Hash2 := sha256.Sum256(sha256Hash1[:])

	// 6. 取前四个字节作为校验和
	checksum := sha256Hash2[:4]

	// 7. 连接版本前缀、哈希结果和校验和
	addressBytes := append(versionPrefixedHash, checksum...)

	// 8. 将字节数组转换为 Base58 编码格式
	address := Base58Encode(addressBytes)

	return address, nil
}

func Base58Encode(input []byte) string {
	var result []byte
	x := new(big.Int).SetBytes(input)

	base := big.NewInt(int64(len(Base58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) > 0 {
		x.DivMod(x, base, mod)
		result = append(result, Base58Alphabet[mod.Int64()])
	}

	// 反转结果
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)

}
