package BLC

func (cli *Rwq_CLI) rwq_send(from []string, to []string, amount []string,nodeID string, mineNow bool) {
	bc := Rwq_NewBlockchain(nodeID)
	defer bc.rwq_db.Close()
	bc.MineNewBlock(from, to, amount,nodeID, mineNow)
}

func (cli *Rwq_CLI) rwq_send_single(from string, to string, amount string,nodeID string, mineNow bool) {
	fromArr := []string{from}
	toArr := []string{to}
	amountArr := []string{amount}

	bc := Rwq_NewBlockchain(nodeID)
	defer bc.rwq_db.Close()
	bc.MineNewBlock(fromArr,toArr,amountArr,nodeID,mineNow)
}