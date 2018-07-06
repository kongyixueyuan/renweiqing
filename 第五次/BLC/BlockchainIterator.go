package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

//迭代器对象
type BlockchainIterator struct {
	CurrentHash []byte
	DB *bolt.DB
}

// 迭代
func (blc *BlockchainIterator) Next() *Block  {
	var block *Block
	err := blc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			//获取当前区块的字节数组
			blockBytes := b.Get(blc.CurrentHash)
			//反序列化
			block = DeserializeBlock(blockBytes)
			//重置迭代器中的hash为上一个区块的hash
			blc.CurrentHash = block.PrevBlockHash
		}
		return nil
	})
	if err!=nil{
		log.Panic(err)
	}
	return block
}
