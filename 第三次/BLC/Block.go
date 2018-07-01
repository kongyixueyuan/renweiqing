package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	//1. 区块高度
	Height int64
	//2. 上一个区块HASH
	PrevBlockHash []byte
	//3. 交易数据
	Data []byte
	//4. 时间戳
	Timestamp int64
	//5. Hash
	Hash []byte
	// 6. Nonce
	Nonce int64
}

//1.创建新的区块
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

// 2.单独写一个方法，生成 创世区块
func CreateGenesisBlock(data string) *Block {
	return NewBlock(data,1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})
}
/*
 * 将区块序列化成字节数组
 */
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

/*
 * 反序列化
 */
func DeserializeBlock(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewBuffer(blockBytes))
	err := decoder.Decode(&block)
	if err != nil{
		log.Panic(err)
	}
	return &block
}