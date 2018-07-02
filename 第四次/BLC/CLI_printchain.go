package BLC


// 终端打印区块
func (cli *CLI) printChain() {
	blockchain := GetBlockchainObject()
	defer blockchain.DB.Close()
	blockchain.Printchain()
}