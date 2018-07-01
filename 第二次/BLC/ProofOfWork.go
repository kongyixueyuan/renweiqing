package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

// 难度：
// 4的倍数  2^4 = 16
// 16进制 一位相当于 二进制的4位
const targetBit = 16

/*
 * 指定区块的工作量证明
 */
type ProofOfWork struct {
	Block  *Block
	target *big.Int
}

//将数据全部转为字节数组后拼接，返回
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.Data,
			IntToHex(pow.Block.Timestamp),
			IntToHex(nonce),
			IntToHex(pow.Block.Height),
		},
		[]byte{},
	)
	return data
}

/*
 * 检查区块中的hash值是否有效
 * proofOfWork.Block.Hash 跟 proofOfWork.Block.target 比较
 */
func (pow *ProofOfWork) IsValid() bool {
	var hashInt big.Int
	hashInt.SetBytes(pow.Block.Hash)
	return pow.target.Cmp(&hashInt) == 1
}

//运行，核心方法，返回，符合要求的hash和nonce
func (pow *ProofOfWork) Run() ([]byte, int64) {
	var nonce int64 = 0
	var hashInt big.Int
	var hash [32]byte
	fmt.Println("开始挖矿中。。。")
	for {
		//将数据全部转为字节数组后拼接
		dataBytes := pow.prepareData(nonce)
		//将拼接后的数据，进行hash运算
		hash = sha256.Sum256(dataBytes)

		//\r是将当前位置移到本行的开头；
		//下面的打印，如果不加\r，会死机
		fmt.Printf("\r%x", hash) //打印 生成 的hash

		// hash转为 big.Int
		hashInt.SetBytes(hash[:]);

		//比较 hashInt 和 target的大小
		// =1 证明 target>hashInt 成立
		if pow.target.Cmp(&hashInt) == 1 {
			break
		}
		// 如果没有匹配的，nonce+1 继续
		nonce = nonce + 1
	}
	fmt.Println("\n挖矿结束。。。")
	return hash[:], nonce
}

/*
 * 创建一个新的工作量证明
*/
func NewProofOfWork(block *Block) *ProofOfWork {
	// 初始化 一个target
	target := big.NewInt(1)

	// 左移 256 - targetBit
	// Lsh sets z = x << n and returns z.
	//func (z *Int) Lsh(x *Int, n uint) *Int
	target = target.Lsh(target, 256-targetBit)

	return &ProofOfWork{block, target}
}
