package main

import (
	"./BLC"
)

func main() {
	//创世区块
	blockchain := BLC.CreateBlockchainWithGenesisBlock()

	cli := &BLC.CLI{blockchain}

	cli.Run()

}
