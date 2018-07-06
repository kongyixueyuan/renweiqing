package BLC

import "bytes"

type TXInput struct {
	TxHash    []byte //  交易的Hash
	Vout      int    //存储TXOutput在Vout里面的索引
	Signature []byte // 数字签名
	PubKey    []byte //公钥，钱包里面
}

// 判断当前的消费是不是自己的钱
func (txInput *TXInput) UnLockRipemd160Hash(ripemd160Hash []byte) bool {

	publicKey := Ripemd160Hash(txInput.PubKey)

	return bytes.Compare(publicKey,ripemd160Hash) == 0
}
