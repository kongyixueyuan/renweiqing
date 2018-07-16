package BLC

import "bytes"

type Rwq_TXOutput struct {
	Rwq_Value  int
	Rwq_PubKeyHash []byte
}
// 根据地址获取 PubKeyHash
func (out *Rwq_TXOutput) Rwq_Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.Rwq_PubKeyHash = pubKeyHash
}

// 判断是否当前公钥对应的交易输出(是否是某个人的交易输出)
func (out *Rwq_TXOutput) Rwq_IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.Rwq_PubKeyHash, pubKeyHash) == 0
}

func Rwq_NewTXOutput(value int, address string) *Rwq_TXOutput {
	txo := &Rwq_TXOutput{value, nil}
	txo.Rwq_Lock([]byte(address))
	return txo
}


