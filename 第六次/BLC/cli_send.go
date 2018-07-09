package BLC

func (cli *Rwq_CLI) rwq_send(from []string, to []string, amount []string) {
	bc := Rwq_NewBlockchain()
	defer bc.rwq_db.Close()
	bc.MineNewBlock(from, to, amount)
}
