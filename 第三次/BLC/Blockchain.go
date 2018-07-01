package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"fmt"
	"time"
	"os"
)

//数据库的名字
const dbName = "blockchain.db"

//表的名字
const blockTableName = "blocks"

type Blockchain struct {
	Tip []byte   //最新区块的Hash
	DB  *bolt.DB //数据库
}

//区块链的迭代器
func (blc *Blockchain) Iterator() *BlockchainIterator  {
	return &BlockchainIterator{blc.Tip,blc.DB}
}
// 判断数据库是否存在
func dbExists() bool  {
	// Stat returns a FileInfo describing the named file.
	// If there is an error, it will be of type *PathError.
	//func Stat(name string) (FileInfo, error)

	//func IsNotExist(err error) bool　　　　　
	//返回一个布尔值，它指明err错误是否报告了一个文件或者目录不存在。它被ErrNotExist 和其它系统调用满足。
	if _,err := os.Stat(dbName);os.IsNotExist(err){
		return false
	}
	return true;
}


// 遍历输出所有区块的信息
func (blc *Blockchain) Printchain()  {
	// 生成迭代器
	blcIterator := blc.Iterator()
	for{
		block := blcIterator.Next()
		fmt.Printf("Height:%d\n",block.Height)
		fmt.Printf("PrevBlockHash:%x\n",block.PrevBlockHash)
		fmt.Printf("Data:%s\n",block.Data)
		fmt.Printf("Timestamp:%s\n",time.Unix( block.Timestamp,0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash:%x\n",block.Hash)
		fmt.Printf("Nonce:%d\n",block.Nonce)
		fmt.Println()

		// 判断 是否到创世区块，退出循环
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}

}

// 增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(data string) {

	err := blc.DB.Update(func(tx *bolt.Tx) error {
		// 1. 获取表
		b := tx.Bucket([]byte(blockTableName))

		// 2. 创建新区块
		if b != nil {
			// 从数据库中取到上一个区块的信息（获取最新区块）
			blockBytes := b.Get(blc.Tip)
			// 反序列化
			block := DeserializeBlock(blockBytes)
			// 3. 将区块序列化，存储到数据库中
			newBlock := NewBlock(data, block.Height+1, block.Hash)
			// 保存生成新的区块到数据库中
			err := b.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			// 4. 更新数据库里面 "l" 对应的hash
			err = b.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			// 5. 更新blockchain的Tip
			blc.Tip = newBlock.Hash
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

}

// 1. 创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *Blockchain {

	if dbExists(){
		fmt.Println("创世区块已经存在，拿到区块链对象(包含最新的区块hash)")
		// 创建或者打开数据库
		db, err := bolt.Open(dbName, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		var blockChain *Blockchain
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blockTableName))
			if b != nil {
				hash := b.Get([]byte("l"))
				blockChain = &Blockchain{hash, db}
			}
			return nil
		})

		//返回区块链对象
		return blockChain
	}
	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var blockHash []byte
	//创建表
	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blockTableName))
		if b == nil{
			// 创建 BlockBucket表
			b, err = tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				log.Panic(err)
			}
		}
		if b != nil {
			//创建创世区块
			genesisBlock := CreateGenesisBlock("Genesis Data......")
			// 将创世区块存储到表中
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//存储最新的区块的hash
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			blockHash = genesisBlock.Hash

		}
		return nil
	})

	//返回区块链对象
	return &Blockchain{blockHash, db}
}
