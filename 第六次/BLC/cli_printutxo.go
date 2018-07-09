package BLC

func (cli *Rwq_CLI) rwq_printutxo() {
	bc := Rwq_NewBlockchain()
	UTXOSet := Rwq_UTXOSet{bc}
	defer bc.rwq_db.Close()
	UTXOSet.String()
}
