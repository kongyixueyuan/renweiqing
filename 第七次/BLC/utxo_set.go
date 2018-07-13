package BLC

import (
	"github.com/boltdb/bolt"
	"log"
	"encoding/hex"
	"fmt"
	"strings"
)

const utxoBucket = "chainstate"

type Rwq_UTXOSet struct {
	Rwq_Blockchain *Rwq_Blockchain
}

// 查询可花费的交易输出
func (u Rwq_UTXOSet) Rwq_FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Rwq_Blockchain.rwq_db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := Rwq_DeserializeOutputs(v)

			for outIdx, out := range outs.Rwq_Outputs {
				if out.Rwq_IsLockedWithKey(pubkeyHash) && accumulated < amount {
					accumulated += out.Rwq_Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return accumulated, unspentOutputs
}

func (u Rwq_UTXOSet) Rwq_Reindex() {
	db := u.Rwq_Blockchain.rwq_db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		// 删除旧的bucket
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic()
		}
		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	UTXO := u.Rwq_Blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(key, outs.Rwq_Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
}

// 生成新区块的时候，更新UTXO数据库
func (u Rwq_UTXOSet) Update(block *Rwq_Block) {
	err := u.Rwq_Blockchain.rwq_db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Rwq_Transactions {
			if !tx.Rwq_IsCoinbase() {
				for _, vin := range tx.Rwq_Vin {
					updatedOuts := Rwq_TXOutputs{}
					outsBytes := b.Get(vin.Rwq_Txid)
					outs := Rwq_DeserializeOutputs(outsBytes)

					// 找出Vin对应的outputs,过滤掉花费的
					for outIndex, out := range outs.Rwq_Outputs {
						if outIndex != vin.Rwq_Vout {
							updatedOuts.Rwq_Outputs = append(updatedOuts.Rwq_Outputs, out)
						}
					}
					// 未花费的交易输出TXOutput为0
					if len(updatedOuts.Rwq_Outputs) == 0 {
						err := b.Delete(vin.Rwq_Txid)
						if err != nil {
							log.Panic(err)
						}
					} else { // 未花费的交易输出TXOutput>0
						err := b.Put(vin.Rwq_Txid, updatedOuts.Rwq_Serialize())
						if err != nil {
							log.Panic(err)
						}
					}
				}
			}

			// 将所有的交易输出TXOutput存入数据库中
			newOutputs := Rwq_TXOutputs{}
			for _, out := range tx.Rwq_Vout {
				newOutputs.Rwq_Outputs = append(newOutputs.Rwq_Outputs, out)
			}
			err := b.Put(tx.Rwq_ID, newOutputs.Rwq_Serialize())
			if err != nil {
				log.Panic(err)
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// 打出某个公钥hash对应的所有未花费输出
func (u *Rwq_UTXOSet) Rwq_FindUTXO(pubKeyHash []byte) []Rwq_TXOutput {
	var UTXOs []Rwq_TXOutput

	err := u.Rwq_Blockchain.rwq_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := Rwq_DeserializeOutputs(v)

			for _, out := range outs.Rwq_Outputs {
				if out.Rwq_IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// 查询某个地址的余额
func (u *Rwq_UTXOSet) Rwq_GetBalance(address string) int {
	balance := 0
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	UTXOs := u.Rwq_FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Rwq_Value
	}
	return balance
}

// 打印所有的UTXO
func (u *Rwq_UTXOSet) String() {
	//outputs := make(map[string][]Rwq_TXOutput)

	var lines []string
	lines = append(lines, "---ALL UTXO:")
	err := u.Rwq_Blockchain.rwq_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := Rwq_DeserializeOutputs(v)

			lines = append(lines, fmt.Sprintf("     Key: %s", txID))
			for i, out := range outs.Rwq_Outputs {
				//outputs[txID] = append(outputs[txID], out)
				lines = append(lines, fmt.Sprintf("     Output: %d", i))
				lines = append(lines, fmt.Sprintf("         value:  %d", out.Rwq_Value))
				lines = append(lines, fmt.Sprintf("         PubKeyHash:  %x", out.Rwq_PubKeyHash))
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(strings.Join(lines, "\n"))
}
