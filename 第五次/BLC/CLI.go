package BLC

import (
	"fmt"
	"os"
	"flag"
	"log"
)

// 终端对象
type CLI struct {
}

// 打印帮助文档
func printUsage() {

	fmt.Println("Usage:")
	fmt.Println("\taddresslists --输出所有钱包地址")
	fmt.Println("\tcreatewallet --创建钱包")
	fmt.Println("\tcreateblockchain -address -amount 100 -- 创建创世区块，可设置Token数量.")
	fmt.Println("\tsend -from FROM -to TO -amount AMOUNT -- 交易明细")
	fmt.Println("\tprintchain -- 输出区块信息.")
	fmt.Println("\tgetbalance -adress -- 查询余额.")
}
// 执行终端
func (cli *CLI) Run() {
	// 判断命令行参数数量是否正确
	isValidArags()

	// 打印所有钱包地址
	addresslistsCmd := flag.NewFlagSet("addresslists",flag.ExitOnError)
	//====================
	// 创建钱包
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	// 添加创世区块命令
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	// 添加创世区块数据命令
	flagCreateBlockchainAddress := createBlockchainCmd.String("address", "", "创建创世区块的地址")
	flagCreateBlockchainAmount := createBlockchainCmd.Int64("amount", 1000, "创建创世区块Token数量")

	//====================
	// 转账命令
	sendBlockCmd := flag.NewFlagSet("send", flag.ExitOnError)

	flagFrom := sendBlockCmd.String("from", "", "转账源地址")
	flagTo := sendBlockCmd.String("to", "", "转账目的地址")
	flagAmount := sendBlockCmd.String("amount", "", "转账金额")

	//====================
	// 打印区块链命令
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	//====================
	//查询余额
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	flagAddress := getBalanceCmd.String("address", "", "地址")
	//====================

	switch os.Args[1] {
	case "addresslists":
		// 检查打印区块链参数是否正确
		err := addresslistsCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		// 打印区块链
		cli.addressLists()
	case "createblockchain":
		// 检查添加区块参数是否正确
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

		if *flagCreateBlockchainAddress == "" {
			fmt.Println("地址不能为空")
			printUsage()
			os.Exit(1)
		}
		if !IsValidForAddress([]byte(*flagCreateBlockchainAddress)) {
			fmt.Println("地址无效。。")
			printUsage()
			os.Exit(1)
		}
		// 创建创世区块
		cli.CreateGenesisBlockchain(*flagCreateBlockchainAddress,*flagCreateBlockchainAmount)
	case "send":
		// 检查添加区块参数是否正确
		err := sendBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		if *flagFrom == "" || *flagTo == "" || *flagAmount == "" {
			printUsage()
			os.Exit(1)
		}
		// 添加区块

		from := JSONToArray(*flagFrom)
		to := JSONToArray(*flagTo)


		for index,fromAdress := range from{
			if !IsValidForAddress([]byte(fromAdress)) || !IsValidForAddress([]byte(to[index])) {
				fmt.Println("地址无效。。")
				os.Exit(1)
			}
		}

		amount := JSONToArray(*flagAmount)

		cli.send(from, to, amount)
	case "printchain":
		// 检查打印区块链参数是否正确
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		// 打印区块链
		cli.printChain()
	case "getbalance":
		// 检查打印区块链参数是否正确
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		if *flagAddress == ""{
			fmt.Println("地址不能为空")
			printUsage()
			os.Exit(1)
		}
		if !IsValidForAddress([]byte(*flagAddress)) {
			fmt.Println("地址无效。。")
			printUsage()
			os.Exit(1)
		}
		// 打印区块链
		cli.getBalance(*flagAddress)
	case "createwallet":
		// 检查打印区块链参数是否正确
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
		if createWalletCmd.Parsed() {
			cli.createWallet()
		}
	default:
		printUsage()
		os.Exit(1)
	}
}



// 判断命令行参数数量是否正确
func isValidArags() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
}
