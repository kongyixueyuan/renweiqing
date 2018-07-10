package BLC

import (
	"github.com/boltdb/bolt"
	"os"
	"fmt"
	"log"
	"encoding/hex"
	"strconv"
	"crypto/ecdsa"
	"bytes"
	"github.com/pkg/errors"
)

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "genesis data 08/07/2018 by viky"

type Rwq_Blockchain struct {
	rwq_tip []byte
	rwq_db  *bolt.DB
}

// 打印区块链内容
func (bc *Rwq_Blockchain) Rwq_Printchain() {
	bci := bc.Rwq_Iterator()

	for {
		block := bci.Rwq_Next()
		block.String()
		if len(block.Rwq_PrevBlockHash) == 0 {
			break
		}
	}

}

// 通过交易hash,查找交易
func (bc *Rwq_Blockchain) Rwq_FindTransaction(ID []byte) (Rwq_Transaction, error) {
	bci := bc.Rwq_Iterator()
	for {
		block := bci.Rwq_Next()
		for _, tx := range block.Rwq_Transactions {
			if bytes.Compare(tx.Rwq_ID, ID) == 0 {
				return *tx, nil
			}
		}
		if len(block.Rwq_PrevBlockHash) == 0 {
			break
		}
	}
	fmt.Printf("查找%x的交易失败",ID)
	return Rwq_Transaction{}, errors.New("未找到交易")
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *Rwq_Blockchain) FindUTXO() map[string]Rwq_TXOutputs {
	// 未花费的交易输出
	// key:交易hash   txID
	UTXO := make(map[string]Rwq_TXOutputs)
	// 已经花费的交易txID : TXOutputs.index
	spentTXOs := make(map[string][]int)
	bci := bc.Rwq_Iterator()

	for {
		block := bci.Rwq_Next()

		// 循环区块中的交易
		for _, tx := range block.Rwq_Transactions {
			// 将区块中的交易hash，转为字符串
			txID := hex.EncodeToString(tx.Rwq_ID)

		Outputs:
			for outIdx, out := range tx.Rwq_Vout { // 循环交易中的 TXOutputs
				// Was the output spent?
				// 如果已经花费的交易输出中，有此输出，证明已经花费
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx { // 如果花费的正好是此笔输出
							continue Outputs // 继续下一次循环
						}
					}
				}

				outs := UTXO[txID] // 获取UTXO指定txID对应的TXOutputs
				outs.Rwq_Outputs = append(outs.Rwq_Outputs, out)
				UTXO[txID] = outs
			}

			if tx.Rwq_IsCoinbase() == false { // 非创世区块
				for _, in := range tx.Rwq_Vin {
					inTxID := hex.EncodeToString(in.Rwq_Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Rwq_Vout)
				}
			}
		}
		// 如果上一区块的hash为0，代表已经到创世区块，循环结束
		if len(block.Rwq_PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// 获取迭代器
func (bc *Rwq_Blockchain) Rwq_Iterator() *Rwq_BlockchainIterator {
	return &Rwq_BlockchainIterator{bc.rwq_tip, bc.rwq_db}
}

// 新建区块链(包含创世区块)
func Rwq_CreateBlockchain(address string) *Rwq_Blockchain {
	if rwq_dbExists(dbFile) {
		fmt.Println("blockchain数据库已经存在.")
		os.Exit(1)
	}

	var tip []byte
	cbtx := Rwq_NewCoinbaseTX(address, genesisCoinbaseData)
	genesis := Rwq_NewGenesisBlock(cbtx)

	//genesis.String()

	// 打开数据库，如果不存在自动创建
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		// 新区块存入数据库
		err = b.Put(genesis.Rwq_Hash, genesis.Rwq_Serialize())
		if err != nil {
			log.Panic(err)
		}
		// 将创世区块的hash存入数据库
		err = b.Put([]byte("l"), genesis.Rwq_Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Rwq_Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Rwq_Blockchain{tip, db}
}

// 获取blockchain对象
func Rwq_NewBlockchain() *Rwq_Blockchain {
	if !rwq_dbExists(dbFile) {
		log.Panic("区块链还未创建")
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Rwq_Blockchain{tip, db}
}

// 生成新的区块(挖矿)
func (bc *Rwq_Blockchain) MineNewBlock(from []string, to []string, amount []string) *Rwq_Block {
	UTXOSet := Rwq_UTXOSet{bc}

	wallets, err := Rwq_NewWallets()
	if err != nil {
		log.Panic(err)
	}

	var txs []*Rwq_Transaction

	for index, address := range from {
		value, _ := strconv.Atoi(amount[index])
		if value<=0 {
			log.Panic("错误：转账金额需要大于0")
		}
		wallet := wallets.Rwq_GetWallet(address)
		tx := Rwq_NewUTXOTransaction(&wallet, to[index], value, &UTXOSet, txs)
		txs = append(txs, tx)
	}

	// 挖矿奖励
	tx := Rwq_NewCoinbaseTX(from[0], "")
	txs = append(txs, tx)

	//=====================================
	var lashHash []byte
	var lastHeight int

	// 检查交易是否有效，验证签名
	for _, tx := range txs {
		if !bc.Rwq_VerifyTransaction(tx,txs) {
			log.Panic("错误：无效的交易")
		}
	}
	// 获取最后一个区块的hash,然后获取最后一个区块的信息，进而获得height
	err = bc.rwq_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lashHash = b.Get([]byte("l"))
		blockData := b.Get(lashHash)
		block := Rwq_DeserializeBlock(blockData)
		lastHeight = block.Rwq_Height
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	// 生成新的区块
	newBlock := Rwq_NewBlock(txs, lashHash, lastHeight+1)

	// 将新区块的内容更新到数据库中
	err = bc.rwq_db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Rwq_Hash,newBlock.Rwq_Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"),newBlock.Rwq_Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.rwq_tip = newBlock.Rwq_Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
	UTXOSet.Update(newBlock)
	return newBlock

}

// 签名
func (bc *Rwq_Blockchain) Rwq_SignTransaction(tx *Rwq_Transaction, privKey ecdsa.PrivateKey,txs []*Rwq_Transaction) {
	prevTXs := make(map[string]Rwq_Transaction)

	// 找到交易输入中，之前的交易
	Vin:
	for _, vin := range tx.Rwq_Vin {
		for _, tx := range txs {
			if bytes.Compare(tx.Rwq_ID, vin.Rwq_Txid) == 0 {
				prevTX := *tx
				prevTXs[hex.EncodeToString(prevTX.Rwq_ID)] = prevTX
				continue Vin
			}
		}

		prevTX, err := bc.Rwq_FindTransaction(vin.Rwq_Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Rwq_ID)] = prevTX

	}

	tx.Rwq_Sign(privKey, prevTXs)
}

// 验证签名
func (bc *Rwq_Blockchain) Rwq_VerifyTransaction(tx *Rwq_Transaction,txs []*Rwq_Transaction) bool {
	if tx.Rwq_IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Rwq_Transaction)
	Vin:
	for _, vin := range tx.Rwq_Vin {
		for _, tx := range txs {
			if bytes.Compare(tx.Rwq_ID, vin.Rwq_Txid) == 0 {
				prevTX := *tx
				prevTXs[hex.EncodeToString(prevTX.Rwq_ID)] = prevTX
				continue Vin
			}
		}
		prevTX, err := bc.Rwq_FindTransaction(vin.Rwq_Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.Rwq_ID)] = prevTX
	}

	return tx.Rwq_Verify(prevTXs)
}

// 判断数据库是否存在
func rwq_dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}
