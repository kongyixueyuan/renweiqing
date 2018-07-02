package BLC


// 创建创世区块
func (cli *CLI) CreateGenesisBlockchain(address string,amount int64) {
	CreateBlockchainWithGenesisBlock(address,amount)
}
