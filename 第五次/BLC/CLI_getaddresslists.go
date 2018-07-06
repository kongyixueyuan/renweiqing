package BLC

import "fmt"

func (cli *CLI) addressLists()  {

	fmt.Println("打印所有的钱包地址:")


	wallets,_ := NewWallets()

	for address,_ := range wallets.WalletsMap{
		fmt.Println(address)
	}
}
