package BLC

import "fmt"

func (cli *Rwq_CLI) rwq_reindexUTXO()  {
	bc := Rwq_NewBlockchain();
	defer bc.rwq_db.Close()
	utxoset := Rwq_UTXOSet{bc}
	utxoset.Rwq_Reindex()
	fmt.Println("重建成功")
}
