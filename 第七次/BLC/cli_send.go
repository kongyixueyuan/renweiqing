package BLC

func (cli *Rwq_CLI) rwq_send(from []string, to []string, amount []string,nodeID string, mineNow bool) {
	bc := Rwq_NewBlockchain(nodeID)
	defer bc.rwq_db.Close()
	bc.MineNewBlock(from, to, amount,nodeID, mineNow)
}