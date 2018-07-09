package BLC

import "log"

func (cli *Rwq_CLI) rwq_createblockchain(address string)  {
	//验证地址是否有效
	if !Rwq_ValidateAddress(address){
		log.Panic("地址无效")
	}
	bc := Rwq_CreateBlockchain(address)
	defer bc.rwq_db.Close()

	// 生成UTXOSet数据库
	UTXOSet := Rwq_UTXOSet{bc}
	UTXOSet.Rwq_Reindex()
}
