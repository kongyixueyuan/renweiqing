package BLC

import "fmt"

type Blockchain struct {
	Blocks []*Block //存储有序的区块
}

/*
 * 生成新的区块并添加到区块链中
 */
func (blc *Blockchain) AddBlockToBlockchain(data string,height int64,preHash []byte)  {
	fmt.Println()
	// 生成新的区块
	block := NewBlock(data,height,preHash)
	// 添加到区块链中
	blc.Blocks = append(blc.Blocks,block)
}

func CreateBlockchainWithGenesisBlock() *Blockchain {
	//创建创世区块
	block := CreateGenesisBlock("Genesis Data ....")
	return &Blockchain{[]*Block{block}}
}