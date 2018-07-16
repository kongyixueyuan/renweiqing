package BLC

import "fmt"

func (cli *Rwq_CLI) rwq_createWallet(nodeID string) {
	//wallet := Rwq_NewWallet()
	//address := wallet.Rwq_GetAddress()
	//fmt.Printf("钱包地址：%s\n",address)

	wallets, _ := Rwq_NewWallets(nodeID)
	address := wallets.Rwq_NewWallet().Rwq_GetAddress()
	wallets.Rwq_SaveToFile(nodeID)
	fmt.Printf("钱包地址：%s\n", address)

}
