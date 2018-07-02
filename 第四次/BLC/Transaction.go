package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

// UTXO
type Transaction struct {

	// 1. 交易hash
	TxHash []byte

	// 2. 输入 记录哪些TXOutput被花费了
	Vins []*TXInput

	// 3. 输出
	Vouts []*TXOutput
}

// 判断当前的交易是否是Coinbase交易 //是否是创世区块里的交易
func (tx *Transaction) IsCoinbaseTransaction() bool {
	return len(tx.Vins[0].TxHash) == 0 && tx.Vins[0].Vout == -1
}

// 1. Transaction 创建分两种情况
// 1. 创世区块创建的Transaction
func NewCoinbaseTransaction(address string ,amount int64) *Transaction {
	// 消费
	txInput := &TXInput{[]byte{},-1,"Genesis Data"}
	txOutput := &TXOutput{amount,address}
	txCoinbase := &Transaction{[]byte{},[]*TXInput{txInput},[]*TXOutput{txOutput}}
	txCoinbase.HashTransaction()
	return txCoinbase
}

// 2. 转账时产生的Transaction
func NewSimpleTransaction(from string,to string,amount int,blockchain *Blockchain,txs []*Transaction) *Transaction {
	// 通过一个函数，返回可用于交易的TXOutput
	money,spendableUTXODic := blockchain.FindSpendableUTXOS(from,amount,txs)

	var txIntputs []*TXInput
	var txOutputs []*TXOutput

	// 获取所有的TXInput
	for txHash,indexArray := range spendableUTXODic  {
		txHashBytes,_ := hex.DecodeString(txHash)
		for _,index := range indexArray  {
			txInput := &TXInput{txHashBytes,index,from}
			txIntputs = append(txIntputs,txInput)
		}
	}

	// 转账
	txOutput := &TXOutput{int64(amount),to}
	txOutputs = append(txOutputs,txOutput)

	// 找零
	txOutput = &TXOutput{int64(money) - int64(amount),from}
	txOutputs = append(txOutputs,txOutput)

	tx := &Transaction{[]byte{},txIntputs,txOutputs}

	//设置hash值
	tx.HashTransaction()

	return tx
}

/*
 * 序列化成字节数组
 */
func (tx *Transaction) HashTransaction() {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(result.Bytes())
	tx.TxHash = hash[:]
}
