package BLC

import (
	"log"
	"fmt"
)

func (cli *Rwq_CLI) rwq_getBalance(address string,nodeID string) {
	if !Rwq_ValidateAddress(address) {
		log.Panic("错误：地址无效")
	}

	bc := Rwq_NewBlockchain(nodeID)
	defer bc.rwq_db.Close()
	UTXOSet := Rwq_UTXOSet{bc}

	balance := UTXOSet.Rwq_GetBalance(address)
	fmt.Printf("地址:%s的余额为：%d\n", address, balance)
}

func (cli *Rwq_CLI) rwq_getBalanceAll(nodeID string) {
	wallets,err := Rwq_NewWallets(nodeID)
	if err!=nil{
		log.Panic(err)
	}
	balances := wallets.Rwq_GetBalanceAll(nodeID)
	for address,balance := range balances{
		fmt.Printf("地址:%s的余额为：%d\n", address, balance)
	}
}