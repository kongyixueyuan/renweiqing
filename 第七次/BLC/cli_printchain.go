package BLC

func (cli *Rwq_CLI) rwq_printchain(nodeID string)  {
	Rwq_NewBlockchain(nodeID).Rwq_Printchain()
}
