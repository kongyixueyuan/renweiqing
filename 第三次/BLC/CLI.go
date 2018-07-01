package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

// 终端对象
type CLI struct {
	BlockChain *Blockchain
}
// 终端添加区块
func (cli *CLI) addBlock(data string)  {
	cli.BlockChain.AddBlockToBlockchain(data)
}
// 终端打印区块
func (cli *CLI) printChain()  {
	cli.BlockChain.Printchain()
}
// 执行终端
func (cli *CLI) Run()  {
	// 判断命令行参数数量是否正确
	isValidArags()
	// 添加区块命令
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	// 打印区块链命令
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	// 添加区块数据命令
	flagAddBlockData := addBlockCmd.String("data", "", "请输入添加到区块的交易数据")

	switch os.Args[1] {
	case "addblock":
		// 检查添加区块参数是否正确
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		if *flagAddBlockData == "" {
			printUsage()
			os.Exit(1)
		}
		// 添加区块
		cli.addBlock(*flagAddBlockData)
	case "printchain":
		// 检查打印区块链参数是否正确
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		// 打印区块链
		cli.printChain()
	default:
		printUsage()
		os.Exit(1)
	}
}

// 打印帮助文档
func printUsage()  {

	fmt.Println("Usage:")
	fmt.Println("\taddblock -data DATA -- 交易数据.")
	fmt.Println("\tprintchain -- 输出区块信息.")
}
// 判断命令行参数数量是否正确
func isValidArags(){
	if len(os.Args) <2 {
		printUsage()
		os.Exit(1)
	}
}