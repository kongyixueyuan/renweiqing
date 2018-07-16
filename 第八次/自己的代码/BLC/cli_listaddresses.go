package BLC

import (
	"log"
	"fmt"
)

func (cli *Rwq_CLI) rwq_listAddrsss(nodeID string)  {
	wallets,err := Rwq_NewWallets(nodeID)

	if err!=nil{
		log.Panic(err)
	}
	addresses := wallets.Rwq_GetAddresses()

	for _,address := range addresses{
		fmt.Printf("%s\n",address)
	}
}
