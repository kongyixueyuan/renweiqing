package BLC

// 转帐
func (cli *CLI) send(from []string, to []string, amount []string) {
	blockchain := GetBlockchainObject()
	defer blockchain.DB.Close()
	blockchain.MineNewBlock(from, to, amount)
}