package BLC

import (
	"math/big"
	"math"
	"bytes"
	"crypto/sha256"
	"fmt"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 16

type Rwq_ProofOfWork struct {
	rwq_block  *Rwq_Block
	rwq_target *big.Int
}

// 生成新的工作量证明
func Rwq_NewProofOfWork(b *Rwq_Block) *Rwq_ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &Rwq_ProofOfWork{b, target}
	return pow
}

// 准备挖矿hash数据
func (pow *Rwq_ProofOfWork) Rwq_PrepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.rwq_block.Rwq_PrevBlockHash,
		pow.rwq_block.Rwq_HashTransactions(),
		IntToHex(pow.rwq_block.Rwq_TimeStamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	}, []byte{})
	return data
}

// 执行工作量证明，返回nonce值和hash
func (pow *Rwq_ProofOfWork) Rwq_Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte

	nonce := 0
	for nonce < maxNonce {
		data := pow.Rwq_PrepareData(nonce)

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		//if math.Remainder(float64(nonce),100000) == 0{
		//	fmt.Printf("\r%x",hash)
		//}
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.rwq_target) == -1 {
			break;
		} else {
			nonce++
		}
	}
	return nonce, hash[:]

}
