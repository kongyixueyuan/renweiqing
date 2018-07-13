package BLC

func (cli *Rwq_CLI) rwq_printutxo(nodeID string) {
	bc := Rwq_NewBlockchain(nodeID)
	UTXOSet := Rwq_UTXOSet{bc}
	defer bc.rwq_db.Close()
	UTXOSet.String()
}
