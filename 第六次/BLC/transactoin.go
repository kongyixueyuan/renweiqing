package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
	"strings"
	"encoding/hex"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/elliptic"
	"math/big"
)

// 创世区块，Token数量
const subsidy  = 10

type Rwq_Transaction struct {
	Rwq_ID   []byte
	Rwq_Vin  []Rwq_TXInput
	Rwq_Vout []Rwq_TXOutput
}

// 是否是创世区块交易
func (tx Rwq_Transaction) Rwq_IsCoinbase() bool {
	// Vin 只有一条
	// Vin 第一条数据的Txid 为 0
	// Vin 第一条数据的Vout 为 -1
	return len(tx.Rwq_Vin) == 1 && len(tx.Rwq_Vin[0].Rwq_Txid) == 0 && tx.Rwq_Vin[0].Rwq_Vout == -1
}
// 将交易序列化
func (tx Rwq_Transaction) Rwq_Serialize() []byte  {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)

	if err != nil{
		log.Panic(err)
	}
	return encoded.Bytes()
}

// 将交易进行Hash
func (tx *Rwq_Transaction) Rwq_Hash() []byte  {
	var hash [32]byte

	txCopy := *tx
	txCopy.Rwq_ID = []byte{}

	hash = sha256.Sum256(txCopy.Rwq_Serialize())
	return hash[:]
}
// 新建创世区块的交易
func Rwq_NewCoinbaseTX(to ,data string) *Rwq_Transaction  {
	if data == ""{
		//如果数据为空，可以随机给默认数据,用于挖矿奖励
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}
	txin := Rwq_TXInput{[]byte{},-1,nil,[]byte(data)}
	txout := Rwq_NewTXOutput(subsidy,to)

	tx := Rwq_Transaction{nil,[]Rwq_TXInput{txin},[]Rwq_TXOutput{*txout}}
	tx.Rwq_ID = tx.Rwq_Hash()
	return &tx
}

// 转帐时生成交易
func Rwq_NewUTXOTransaction(wallet *Rwq_Wallet,to string,amount int,UTXOSet *Rwq_UTXOSet,txs []*Rwq_Transaction) *Rwq_Transaction   {

	pubKeyHash := Rwq_HashPubKey(wallet.Rwq_PublicKey)
	acc, validOutputs := UTXOSet.Rwq_FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("账户余额不足")
	}

	var inputs []Rwq_TXInput
	var outputs []Rwq_TXOutput
	// 构造input
	for txid,outs := range validOutputs{
		txID,err := hex.DecodeString(txid)
		if err !=nil{
			log.Panic(err)
		}

		for _,out := range outs{
			input := Rwq_TXInput{txID,out,nil,wallet.Rwq_PublicKey}
			inputs = append(inputs,input)
		}
	}
	// 生成交易输出
	outputs = append(outputs,*Rwq_NewTXOutput(amount,to))
	// 生成余额
	if acc > amount {
		outputs = append(outputs,*Rwq_NewTXOutput(acc-amount,string(wallet.Rwq_GetAddress())))
	}

	tx := Rwq_Transaction{nil,inputs,outputs}
	tx.Rwq_ID = tx.Rwq_Hash()
	// 签名
	UTXOSet.Rwq_Blockchain.Rwq_SignTransaction(&tx,wallet.Rwq_PrivateKey)

	return &tx

}

// 签名
func (tx *Rwq_Transaction) Rwq_Sign(privateKey ecdsa.PrivateKey,prevTXs map[string]Rwq_Transaction)  {
	if tx.Rwq_IsCoinbase() { // 创世区块不需要签名
		return
	}

	// 检查交易的输入是否正确
	for _,vin := range tx.Rwq_Vin{
		if prevTXs[hex.EncodeToString(vin.Rwq_Txid)].Rwq_ID == nil{
			log.Panic("错误：之前的交易不正确")
		}
	}

	txCopy := tx.Rwq_TrimmedCopy()

	for inID, vin := range txCopy.Rwq_Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Rwq_Txid)]
		txCopy.Rwq_Vin[inID].Rwq_Signature = nil
		txCopy.Rwq_Vin[inID].Rwq_PubKey = prevTx.Rwq_Vout[vin.Rwq_Vout].Rwq_PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, []byte(dataToSign))
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Rwq_Vin[inID].Rwq_Signature = signature
		txCopy.Rwq_Vin[inID].Rwq_PubKey = nil
	}
}
// 验证签名
func (tx *Rwq_Transaction) Rwq_Verify(prevTXs map[string]Rwq_Transaction) bool {
	if tx.Rwq_IsCoinbase() {
		return true
	}

	for _, vin := range tx.Rwq_Vin {
		if prevTXs[hex.EncodeToString(vin.Rwq_Txid)].Rwq_ID == nil {
			log.Panic("错误：之前的交易不正确")
		}
	}

	txCopy := tx.Rwq_TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Rwq_Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Rwq_Txid)]
		txCopy.Rwq_Vin[inID].Rwq_Signature = nil
		txCopy.Rwq_Vin[inID].Rwq_PubKey = prevTx.Rwq_Vout[vin.Rwq_Vout].Rwq_PubKeyHash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Rwq_Signature)
		r.SetBytes(vin.Rwq_Signature[:(sigLen / 2)])
		s.SetBytes(vin.Rwq_Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.Rwq_PubKey)
		x.SetBytes(vin.Rwq_PubKey[:(keyLen / 2)])
		y.SetBytes(vin.Rwq_PubKey[(keyLen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.Rwq_Vin[inID].Rwq_PubKey = nil
	}

	return true
}

// 复制交易（输入的签名和公钥置为空）
func (tx *Rwq_Transaction) Rwq_TrimmedCopy() Rwq_Transaction {
	var inputs []Rwq_TXInput
	var outputs []Rwq_TXOutput

	for _, vin := range tx.Rwq_Vin {
		inputs = append(inputs, Rwq_TXInput{vin.Rwq_Txid, vin.Rwq_Vout, nil, nil})
	}

	for _, vout := range tx.Rwq_Vout {
		outputs = append(outputs, Rwq_TXOutput{vout.Rwq_Value, vout.Rwq_PubKeyHash})
	}

	txCopy := Rwq_Transaction{tx.Rwq_ID, inputs, outputs}

	return txCopy
}
// 打印交易内容
func (tx Rwq_Transaction) String()  {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction ID: %x", tx.Rwq_ID))

	for i, input := range tx.Rwq_Vin {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Rwq_Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Rwq_Vout))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Rwq_Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.Rwq_PubKey))
	}

	for i, output := range tx.Rwq_Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Rwq_Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.Rwq_PubKeyHash))
	}
	fmt.Println(strings.Join(lines, "\n"))
}