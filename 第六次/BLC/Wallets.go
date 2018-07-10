package BLC

import (
	"os"
	"io/ioutil"
	"log"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
)

const walletFile  = "wallet.dat"

type Rwq_Wallets struct {
	Rwq_Wallets map[string]*Rwq_Wallet
}

// 生成新的钱包
// 从数据库中读取，如果不存在
func Rwq_NewWallets()(*Rwq_Wallets,error)  {
	wallets := Rwq_Wallets{}
	wallets.Rwq_Wallets = make(map[string]*Rwq_Wallet)

	err := wallets.Rwq_LoadFromFile()

	return &wallets,err
}
// 生成新的钱包地址列表
func (ws *Rwq_Wallets) Rwq_NewWallet() *Rwq_Wallet {
	wallet := Rwq_NewWallet()
	address := wallet.Rwq_GetAddress()
	ws.Rwq_Wallets[string(address)] = wallet
	return wallet
}
// 获取钱包地址
func (ws *Rwq_Wallets) Rwq_GetAddresses()[]string  {
	var addresses []string
	for address := range ws.Rwq_Wallets{
		addresses = append(addresses,address)
	}
	return addresses
}

// 根据地址获取钱包的详细信息
func (ws Rwq_Wallets) Rwq_GetWallet(address string) Rwq_Wallet {
	return *ws.Rwq_Wallets[address]
}

// 从数据库中读取钱包列表
func (ws *Rwq_Wallets)Rwq_LoadFromFile() error  {
	 if _,err := os.Stat(walletFile) ; os.IsNotExist(err){
	 	return err
	 }

	 fileContent ,err := ioutil.ReadFile(walletFile)
	 if err !=nil{
	 	log.Panic(err)
	 }

	 var wallets Rwq_Wallets
	 gob.Register(elliptic.P256())
	 decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	 err = decoder.Decode(&wallets)
	 if err !=nil{
	 	log.Panic(err)
	 }

	 ws.Rwq_Wallets = wallets.Rwq_Wallets

	 return nil
}

// 将钱包存到数据库中
func (ws *Rwq_Wallets)Rwq_SaveToFile()  {
	var content bytes.Buffer

	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err !=nil{
		log.Panic(err)
	}

	err = ioutil.WriteFile(walletFile,content.Bytes(),0644)
	if err !=nil{
		log.Panic(err)
	}
}
// 打印所有钱包的余额
func (ws *Rwq_Wallets) Rwq_GetBalanceAll() map[string]int {
	addresses := ws.Rwq_GetAddresses()
	bc := Rwq_NewBlockchain()
	defer bc.rwq_db.Close()
	UTXOSet := Rwq_UTXOSet{bc}

	result := make(map[string]int)
	for _,address := range addresses{
		if !Rwq_ValidateAddress(address) {
			result[address] = -1
		}
		balance := UTXOSet.Rwq_GetBalance(address)
		result[address] = balance
	}
	return result
}