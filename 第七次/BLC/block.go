package BLC

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
	"fmt"
)

type Rwq_Block struct {
	Rwq_TimeStamp     int64
	Rwq_Transactions   []*Rwq_Transaction
	Rwq_PrevBlockHash []byte
	Rwq_Hash          []byte
	Rwq_Nonce         int
	Rwq_Height        int
}
// 生成新的区块
func Rwq_NewBlock(transactions []*Rwq_Transaction, prevBlockHash []byte, height int) *Rwq_Block {
	// 生成新的区块对象
	block := &Rwq_Block{
		time.Now().Unix(),
		transactions,
		prevBlockHash,
		[]byte{},
		0,
		height,
	}
	// 挖矿

	pow := Rwq_NewProofOfWork(block)
	nonce,hash :=pow.Rwq_Run()

	block.Rwq_Nonce = nonce
	block.Rwq_Hash = hash[:]

	return block

}

// 将交易进行hash
func (b Rwq_Block) Rwq_HashTransactions() []byte {
	var transactions [][]byte
	// 获取交易真实内容
	for _,tx := range b.Rwq_Transactions{
		transactions = append(transactions,tx.Rwq_Serialize())
	}
	//txHash := sha256.Sum256(bytes.Join(transactions,[]byte{}))
	mTree := Rwq_NewMerkelTree(transactions)
	return mTree.Rwq_RootNode.Rwq_Data
}
// 新建创世区块
func Rwq_NewGenesisBlock(coinbase *Rwq_Transaction) *Rwq_Block  {
	return Rwq_NewBlock([]*Rwq_Transaction{coinbase},[]byte{},1)
}

// 序列化
func (b *Rwq_Block) Rwq_Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// 反序列化
func Rwq_DeserializeBlock(d []byte) *Rwq_Block {
	var block Rwq_Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
// 打印区块内容
func (block Rwq_Block) String()  {
	fmt.Println("\n==============")
	fmt.Printf("Height:\t%d\n", block.Rwq_Height)
	fmt.Printf("PrevBlockHash:\t%x\n", block.Rwq_PrevBlockHash)
	fmt.Printf("Timestamp:\t%s\n", time.Unix(block.Rwq_TimeStamp, 0).Format("2006-01-02 03:04:05 PM"))
	fmt.Printf("Hash:\t%x\n", block.Rwq_Hash)
	fmt.Printf("Nonce:\t%d\n", block.Rwq_Nonce)
	fmt.Println("Txs:")

	for _, tx := range block.Rwq_Transactions {
		tx.String()
	}
	fmt.Println("==============")
}
