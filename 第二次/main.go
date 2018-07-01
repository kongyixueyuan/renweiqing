package main

import (
	"./BLC"
	"fmt"
	"time"
)

func main()  {
	//生成带创世区块的区块链
	fmt.Println("创建创世区块中。。。")
	blockchain  := BLC.CreateBlockchainWithGenesisBlock()
	fmt.Println("创世区块创建成功")

	//生成一个新的区块
	blockchain.AddBlockToBlockchain("viky send 100BTC to lulu",blockchain.Blocks[len(blockchain.Blocks)-1].Height+1,blockchain.Blocks[len(blockchain.Blocks)-1].Hash)

	fmt.Println("\n链中总共有" , len(blockchain.Blocks),"个区块")
	fmt.Println("区块中的数据都有：")
	for _,val := range blockchain.Blocks {
		println("第",val.Height,"个区块:")
		println(string(val.Data))
		println("=============================")
	}
	time.Sleep(1*time.Second)
	fmt.Println()
	//增加工作量证明，验证

	newBlock := BLC.NewBlock("测试工作量证明",1,[]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0})

	fmt.Println("\n区块内容如下：")
	fmt.Printf("data:%s\n",newBlock.Data)
	fmt.Printf("nonce:%d\n",newBlock.Nonce)
	fmt.Printf("hash:%x\n\n",newBlock.Hash)
	fmt.Println(newBlock)

	pow := BLC.NewProofOfWork(newBlock)
	if pow.IsValid() {
		fmt.Println("是有效区块")
	}else{
		fmt.Println("不是有效区块")
	}
}
