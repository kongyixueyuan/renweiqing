package BLC

import "fmt"

func (cli *Rwq_CLI) rwq_reindexUTXO(nodeID string)  {
	bc := Rwq_NewBlockchain(nodeID);
	defer bc.rwq_db.Close()
	utxoset := Rwq_UTXOSet{bc}
	utxoset.Rwq_Reindex()
	fmt.Println("重建成功")
}
