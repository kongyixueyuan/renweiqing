package BLC

import "time"

type Block struct {
	// 区块高度  或者 理解为第个区块
	Height int64
	// 上一个区块的hash值
	PrevBlockHash []byte
	//  数据
	Data []byte
	// 时间戳
	Timestamp int64
	// Hash
	Hash []byte
	// nonce : 循环生成Hash用
	Nonce int64
}
//创建新的区块
func NewBlock(data string , height int64, prevBlockHash []byte)  *Block {
	//创建区块
	block := &Block{Height:height,PrevBlockHash:prevBlockHash,Data:[]byte(data),Timestamp:time.Now().Unix(),Hash:nil,Nonce:0}
	// 调用工作证明的方法 生成hash 和 nonce 并更新到区块中
	// 生成一个新的工作证明
	pow := NewProofOfWork(block)
	// 运行工作证明
	hash,nonce :=  pow.Run();

	// 将工作证明的结果 ，赋值给区块
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

/*
 * 创建创世区块
 */
func CreateGenesisBlock(data string) *Block  {
	return NewBlock(data,1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}