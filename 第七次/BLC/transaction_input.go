package BLC

import "bytes"

type Rwq_TXInput struct {
	Rwq_Txid      []byte
	Rwq_Vout      int      // Vout的index
	Rwq_Signature []byte   // 签名
	Rwq_PubKey    []byte   // 公钥
}

func (in Rwq_TXInput) UsesKey(pubKeyHash []byte) bool  {
	lockingHash := Rwq_HashPubKey(in.Rwq_PubKey)

	return bytes.Compare(lockingHash,pubKeyHash) == 0
}
