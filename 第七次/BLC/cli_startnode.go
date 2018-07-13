package BLC

import (
	"fmt"
	"log"
)

func (cli *Rwq_CLI) rwq_startNode(nodeID, minerAddress string) {
	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if Rwq_ValidateAddress(minerAddress) {
			fmt.Println("挖矿开始，挖矿地址为: ", minerAddress)
		} else {
			log.Panic("挖矿地址错误!")
		}
	}
	Rwq_StartServer(nodeID, minerAddress)
}
