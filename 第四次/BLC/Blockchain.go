package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"fmt"
	"time"
	"os"
	"encoding/hex"
	"strconv"
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
func (blc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blc.Tip, blc.DB}
}

// 判断数据库是否存在
func DBExists() bool {
	// Stat returns a FileInfo describing the named file.
	// If there is an error, it will be of type *PathError.
	//func Stat(name string) (FileInfo, error)

	//func IsNotExist(err error) bool　　　　　
	//返回一个布尔值，它指明err错误是否报告了一个文件或者目录不存在。它被ErrNotExist 和其它系统调用满足。
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true;
}

// 遍历输出所有区块的信息
func (blc *Blockchain) Printchain() {
	// 生成迭代器
	blcIterator := blc.Iterator()
	for {
		fmt.Println("==============")
		block := blcIterator.Next()
		fmt.Printf("Height:%d\n", block.Height)
		fmt.Printf("PrevBlockHash:%x\n", block.PrevBlockHash)
		fmt.Printf("Txs:%v\n", block.Txs)
		fmt.Printf("Timestamp:%s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 03:04:05 PM"))
		fmt.Printf("Hash:%x\n", block.Hash)
		fmt.Printf("Nonce:%d\n", block.Nonce)
		fmt.Println("Txs:")

		for _, tx := range block.Txs {
			fmt.Printf("TxHash:%v\n", tx.TxHash)
			fmt.Println("Vins:")
			for _, in := range tx.Vins {
				fmt.Printf("TxHash:%s,Vout:%d,ScriptSig:%s\n", in.TxHash, in.Vout, in.ScriptSig)
			}
			fmt.Println("Vouts:")
			for _, out := range tx.Vouts {
				fmt.Println("Value:", out.Value, "  ScriptPubKey:", out.ScriptPubKey)
			}

		}
		fmt.Println("==============")

		// 判断 是否到创世区块，退出循环
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}

}

// 增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(txs []*Transaction) {

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
			newBlock := NewBlock(txs, block.Height+1, block.Hash)
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
func CreateBlockchainWithGenesisBlock(address string,amount int64) *Blockchain {

	//判断数据库是否存在，如果存在证明已经创建过创世区块，返回
	if DBExists() {
		fmt.Println("创世区块已经存在。。。")
		os.Exit(1)
	}

	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var genesisHash []byte
	//创建表
	err = db.Update(func(tx *bolt.Tx) error {

		fmt.Println("正在创建创世区块")

		// 创建 BlockBucket表
		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}

		if b != nil {
			// 创建一个coinbase Transaction
			txCoinbase := NewCoinbaseTransaction(address,amount)
			//创建创世区块
			genesisBlock := CreateGenesisBlock([]*Transaction{txCoinbase})
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
			genesisHash = genesisBlock.Hash
		}
		return nil
	})
	return &Blockchain{genesisHash, db}
}

func GetBlockchainObject() *Blockchain {
	//如果数据库不存在 ，无法获取区块链
	if ! DBExists() {
		fmt.Println("数据库不存在,创世区块还未生成")
		os.Exit(1)
	}

	// 创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var tip []byte
	//查看表
	err = db.View(func(tx *bolt.Tx) error {
		// 创建 BlockBucket表
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			// 将创世区块存储到表中
			tip = b.Get([]byte("l"))
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	return &Blockchain{tip, db}
}

// 如果一个地址对应的TXOutput未花费，那么这个Transaction就应该添加到数组中返回
// 查找未花费的TXOutput
func (blockchain *Blockchain) UnUTXOs(address string, txs []*Transaction) []*UTXO {

	var unUTXOs []*UTXO

	//花费的 交易输出
	spentTXOutputs := make(map[string][]int)

	for _, tx := range txs {
		// 如果是创世区块，不考虑Vins
		if tx.IsCoinbaseTransaction() == false {
			for _, in := range tx.Vins {
				//是否是指定地址的 TXInput
				if in.UnLockWithAddress(address) {
					// 转字符串
					key := hex.EncodeToString(in.TxHash)
					// 如果是，记录到spentTXOutputs
					spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
				}

			}
		}
	}

	for _, tx := range txs {

	Work1:
		for index, out := range tx.Vouts {
			// 如果是指定地址的 TXOutputs
			if out.UnLockScriptPubKeyWithAddress(address) {
				// 如果花费的记录为0，代表没有花费过，直接返回
				if len(spentTXOutputs) == 0 {
					utxo := &UTXO{tx.TxHash, index, out}
					unUTXOs = append(unUTXOs, utxo)
				} else {
					for hash, indexArray := range spentTXOutputs {
						txHashStr := hex.EncodeToString(tx.TxHash)
						// 如果 花费记录中的hash值与，当前交易的hash一样，说明此交易记录已经被花费
						if hash == txHashStr {
							var isUnSpentUTXO bool

							for _, outIndex := range indexArray {
								// 索引相同，代表已花费
								if index == outIndex {
									isUnSpentUTXO = true
									continue Work1
								}
							}
							if isUnSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}
						} else { //不等于代表未花费，直接添加
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}
					}
				}

			}

		}

	}

	// 遍历所有区块
	// 下面代码同上面的逻辑
	blockIterator := blockchain.Iterator()
	for {
		block := blockIterator.Next()
		for i := len(block.Txs) - 1; i >= 0; i-- {
			tx := block.Txs[i]
			if tx.IsCoinbaseTransaction() == false {
				for _, in := range tx.Vins {
					//是否能够解锁
					if in.UnLockWithAddress(address) {
						key := hex.EncodeToString(in.TxHash)
						spentTXOutputs[key] = append(spentTXOutputs[key], in.Vout)
					}
				}
			}
			// Vouts

		work:
			for index, out := range tx.Vouts {
				if out.UnLockScriptPubKeyWithAddress(address) {
					if spentTXOutputs != nil {
						if len(spentTXOutputs) != 0 {
							var isSpentUTXO bool
							for txHash, indexArray := range spentTXOutputs {
								for _, i := range indexArray {
									if index == i && txHash == hex.EncodeToString(tx.TxHash) {
										isSpentUTXO = true
										continue work
									}
								}
							}
							if isSpentUTXO == false {
								utxo := &UTXO{tx.TxHash, index, out}
								unUTXOs = append(unUTXOs, utxo)
							}
						} else {
							utxo := &UTXO{tx.TxHash, index, out}
							unUTXOs = append(unUTXOs, utxo)
						}

					}
				}

			}

		}

		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		//如果当前区块是创世区块，退出循环
		if hashInt.Cmp(big.NewInt(0)) == 0 {
			break;
		}

	}

	return unUTXOs
}

// 转账时查找可用的UTXO
func (blockchain *Blockchain) FindSpendableUTXOS(from string, amount int, txs []*Transaction) (int64, map[string][]int) {

	//1. 现获取所有的UTXO
	utxos := blockchain.UnUTXOs(from, txs)

	spendableUTXO := make(map[string][]int)

	//2. 遍历utxos
	var value int64
	for _, utxo := range utxos {
		value = value + utxo.Output.Value
		hash := hex.EncodeToString(utxo.TxHash)
		spendableUTXO[hash] = append(spendableUTXO[hash], utxo.Index)
		if value >= int64(amount) {
			break
		}
	}
	//如果最后可用金额小于交易金额，提示余额不足
	if value < int64(amount) {

		fmt.Printf("%s's fund is 不足\n", from)
		os.Exit(1)
	}

	return value, spendableUTXO
}

//挖新的区块
func (blockchain *Blockchain) MineNewBlock(from []string, to []string, amount []string) {

	var txs []*Transaction
	//根据 输入参数，生成交易
	for index, address := range from {
		// Atoi returns the result of ParseInt(s, 10, 0) converted to type int.
		value, _ := strconv.Atoi(amount[index])
		// 生成新的交易
		tx := NewSimpleTransaction(address, to[index], value, blockchain, txs)
		// 打包交易
		txs = append(txs, tx)
	}

	//1. 通过相关算法建立Transaction数组
	var block *Block
	//获取最新的区块
	blockchain.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			hash := b.Get([]byte("l"))
			blockBytes := b.Get(hash)
			block = DeserializeBlock(blockBytes)
		}
		return nil
	})

	//2. 建立新的区块
	block = NewBlock(txs, block.Height+1, block.Hash)

	//将新区块存储到数据库
	blockchain.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockTableName))
		if b != nil {
			b.Put(block.Hash, block.Serialize())
			b.Put([]byte("l"), block.Hash)
			blockchain.Tip = block.Hash
		}
		return nil
	})
}

// 查询余额
func (blockchain *Blockchain) GetBalance(address string) int64 {
	// 获取当前地址的"未花费的金额"
	// 然后将"未花费的金额"相加
	utxos := blockchain.UnUTXOs(address, []*Transaction{})

	var amount int64
	for _, utxo := range utxos {
		amount = amount + utxo.Output.Value
	}
	return amount
}
