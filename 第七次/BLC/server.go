package BLC

import (
	"fmt"
	"net"
	"log"
	"bytes"
	"encoding/gob"
	"io"
	"io/ioutil"
	"encoding/hex"
)

const protocol = "tcp"   // 节点协议
const nodeVersion = 1    // 节点版本
const commandLength = 12 // 命令长度

var nodeAddress string                         // 当前节点地址
var miningAddress string                       // 挖矿地址
var knownNodes = []string{"localhost:3000"}    // 存储所有已知节点
var blocksInTransit = [][]byte{}               // 存储接收到的区块hash
var mempool = make(map[string]Rwq_Transaction) // 存储接收到的交易

type rwq_addr struct {
	Rwq_AddrList []string
}

type rwq_block struct {
	Rwq_AddrFrom string
	Rwq_Block    []byte
}

type rwq_getblocks struct {
	Rwq_AddrFrom string
}

type rwq_getdata struct {
	Rwq_AddrFrom string
	Rwq_Type     string
	Rwq_ID       []byte
}

type rwq_inv struct {
	Rwq_AddrFrom string
	Rwq_Type     string
	Rwq_Items    [][]byte
}

type rwq_txs struct {
	Rwq_AddFrom     string
	Rwq_Transactions [][]byte
}

type rwq_version struct {
	Rwq_Version    int
	Rwq_BestHeight int
	Rwq_AddrFrom   string
}

//启动Server
func Rwq_StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := Rwq_NewBlockchain(nodeID)

	// 3000端口为：主节点
	// 3001端口为：钱包节点
	// 3002端口为：挖矿节点
	if nodeAddress != knownNodes[0] {
		// 此节点是钱包节点或者矿工节点，需要向主节点发送请求同步数据
		rwq_sendVersion(knownNodes[0], bc)
	}

	for { // 等待接收命令
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go rwq_handleConnecton(conn, bc)
	}
}

// 向中心节点发送 version 消息来查询是否自己的区块链已过时
func rwq_sendVersion(addr string, bc *Rwq_Blockchain) {
	bestHeight := bc.Rwq_GetBestHeight()
	payload := rwq_gobEncode(rwq_version{nodeVersion, bestHeight, nodeAddress})

	request := append(rwq_commandToBytes("version"), payload...)

	rwq_sendData(addr, request)
}

// 发送数据
func rwq_sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	// 如果连接失败，更新节点数据
	if err != nil {
		fmt.Sprintf("%s地址不可用\n", addr)
		var updatedNodes []string

		for _, node := range knownNodes {
			if node != addr {
				updatedNodes = append(updatedNodes, node)
			}
		}

		knownNodes = updatedNodes
		return
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}

}

// 发送获取区块的的命令
func rwq_sendGetBlocks(address string) {
	payload := rwq_gobEncode(rwq_getblocks{nodeAddress})
	request := append(rwq_commandToBytes("getblocks"), payload...)

	rwq_sendData(address, request)
}

// 发送处理区块及交易hash列表请求
func rwq_sendInv(address, kind string, items [][]byte) {
	inventory := rwq_inv{nodeAddress, kind, items}
	payload := rwq_gobEncode(inventory)
	request := append(rwq_commandToBytes("inv"), payload...)

	rwq_sendData(address, request)
}

// 发送获取区块详细数据的命令
func rwq_sendGetData(address, kind string, id []byte) {
	payload := rwq_gobEncode(rwq_getdata{nodeAddress, kind, id})
	request := append(rwq_commandToBytes("getdata"), payload...)

	rwq_sendData(address, request)
}

// 发送区块内容命令
func rwq_sendBlock(addr string, b *Rwq_Block) {
	data := rwq_block{nodeAddress, b.Rwq_Serialize()}
	payload := rwq_gobEncode(data)
	request := append(rwq_commandToBytes("block"), payload...)

	rwq_sendData(addr, request)
}

// 发送交易内容命令
func rwq_sendTx(addr string, tx *Rwq_Transaction) {
	txs := []*Rwq_Transaction{tx}
	rwq_sendTxs(addr,txs)
}
// 发送交易内容命令
func rwq_sendTxs(addr string, txs []*Rwq_Transaction) {

	data := rwq_txs{nodeAddress, Rwq_SerializeTransactions(txs)}
	payload := rwq_gobEncode(data)
	request := append(rwq_commandToBytes("tx"), payload...)

	rwq_sendData(addr, request)
}

//================================================================
// 命令收集并分发
func rwq_handleConnecton(conn net.Conn, bc *Rwq_Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := rwq_bytesToCommand(request[:commandLength])
	fmt.Printf("接收到命令：%s", command)

	switch command {
	case "addr": // 添加新节点
		rwq_handleAddr(request)
	case "block": // 添加新区块
		rwq_handleBlock(request, bc)
	case "inv": // 接收区块及交易hash列表 ，区块获取到内容到存储到 blocksInTransit ， 交易存储到 mempool
		rwq_handleInv(request, bc)
	case "getblocks": // 将当前节点区块链中的所有区块hash，返回给请求节点
		rwq_handleGetBlocks(request, bc)
	case "getdata": // 将单个交易或区块的内容 返回给请求节点
		rwq_handleGetData(request, bc)
	case "tx": // 添加新的交易,交易数量大于2，矿工节点挖矿,如果是主节点，进行分发交易
		rwq_handleTx(request, bc)
	case "version": // 检查是否需要同步数据，根据区块的height
		rwq_handleVersion(request, bc)
	default:
		fmt.Println("未知命令!")
	}

	conn.Close()

}

// 添加新节点
func rwq_handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload rwq_addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.Rwq_AddrList...)
	fmt.Printf("有%d个节点加入\n", len(knownNodes))
	// 把新节点推送给其他所有节点
	for _, node := range knownNodes {
		rwq_sendGetBlocks(node)
	}
}

/*
当接收到一个新块时，我们把它放到区块链里面。
如果还有更多的区块需要下载，我们继续从上一个下载的块的那个节点继续请求。
当最后把所有块都下载完后，对 UTXO 集进行重新索引

TODO: 并非无条件信任，我们应该在将每个块加入到区块链之前对它们进行验证。
TODO: 并非运行 UTXOSet.Reindex()， 而是应该使用 UTXOSet.Update(block)，
TODO: 因为如果区块链很大，它将需要很多时间来对整个 UTXO 集重新索引。
 */
func rwq_handleBlock(request []byte, bc *Rwq_Blockchain) {
	var buff bytes.Buffer
	var payload rwq_block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Rwq_Block
	block := Rwq_DeserializeBlock(blockData)

	fmt.Println("收到一个新的区块!")
	bc.Rwq_AddBlock(block)

	fmt.Printf("Added block %x\n", block.Rwq_Hash)

	// 如果还有更多的区块需要下载，继续从上一个下载的块的那个节点继续请求
	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		rwq_sendGetData(payload.Rwq_AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := Rwq_UTXOSet{bc}
		UTXOSet.Rwq_Reindex()
	}
}

// 向其他节点展示当前节点有什么块和交易,接收区块及交易列表
func rwq_handleInv(request []byte, bc *Rwq_Blockchain) {
	var buff bytes.Buffer
	var payload rwq_inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("接收到列表,长度为：%d，类型： %s\n", len(payload.Rwq_Items), payload.Rwq_Type)

	// 如果数据是 区块
	if payload.Rwq_Type == "block" {
		blocksInTransit = payload.Rwq_Items

		blockHash := payload.Rwq_Items[0]
		// 发送获取区块详细数据的命令
		rwq_sendGetData(payload.Rwq_AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}
	// 如果数据是 交易
	// 转账时，未立即挖矿，将交易保存到内存池中
	if payload.Rwq_Type == "tx" {
		txID := payload.Rwq_Items[0]
		// 如果内存池中，交易内容为空
		if mempool[hex.EncodeToString(txID)].Rwq_ID == nil {
			// 发送获取交易详细内容请求
			rwq_sendGetData(payload.Rwq_AddrFrom, "tx", txID)
		}
	}
}

// 处理获取区块命令，返回区块链中的所有区块hash
func rwq_handleGetBlocks(request []byte, bc *Rwq_Blockchain) {
	var buff bytes.Buffer
	var payload rwq_getblocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.Rwq_GetBlockHashes()
	rwq_sendInv(payload.Rwq_AddrFrom, "block", blocks)
}

//  将单个交易或区块的内容 返回给请求节点
func rwq_handleGetData(request []byte, bc *Rwq_Blockchain) {
	var buff bytes.Buffer
	var payload rwq_getdata

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Rwq_Type == "block" {
		block, err := bc.Rwq_GetBlock([]byte(payload.Rwq_ID))
		if err != nil {
			return
		}

		rwq_sendBlock(payload.Rwq_AddrFrom, &block)
	}

	if payload.Rwq_Type == "tx" {
		txID := hex.EncodeToString(payload.Rwq_ID)
		tx := mempool[txID]

		rwq_sendTx(payload.Rwq_AddrFrom, &tx)
		// delete(mempool, txID)
	}
}

// 处理交易
func rwq_handleTx(request []byte, bc *Rwq_Blockchain) {
	var buff bytes.Buffer
	var payload rwq_txs

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Rwq_Transactions
	txsDes := Rwq_DeserializeTransactions(txData)

	for _,tx := range txsDes {
		mempool[hex.EncodeToString(tx.Rwq_ID)] = tx

		// 如果是主节点
		if nodeAddress == knownNodes[0] {
			for _, node := range knownNodes {
				// 给其他节点分发，添加交易
				if node != nodeAddress && node != payload.Rwq_AddFrom {
					rwq_sendInv(node, "tx", [][]byte{tx.Rwq_ID})
				}
			}
		} else {
			// 如果交易池中有两条交易 并且当前是挖矿节点
			if len(mempool) >= 2 && len(miningAddress) > 0 {
			MineTransactions:
				var txs []*Rwq_Transaction

				for id := range mempool {
					tx := mempool[id]
					if bc.Rwq_VerifyTransaction(&tx, txs) {
						txs = append(txs, &tx)
					}
				}

				if len(txs) == 0 {
					fmt.Println("交易不可用...")
					break
				}

				cbTx := Rwq_NewCoinbaseTX(miningAddress, "")
				txs = append(txs, cbTx)

				newBlock := bc.Rwq_MineBlock(txs)
				UTXOSet := Rwq_UTXOSet{bc}
				UTXOSet.Update(newBlock)

				fmt.Println("挖到新的区块!")

				for _, tx := range txs {
					txID := hex.EncodeToString(tx.Rwq_ID)
					delete(mempool, txID)
				}

				for _, node := range knownNodes {
					if node != nodeAddress {
						rwq_sendInv(node, "block", [][]byte{newBlock.Rwq_Hash})
					}
				}

				if len(mempool) > 0 {
					goto MineTransactions
				}
			}
		}
	}
}

// 检查是否需要同步数据
func rwq_handleVersion(request []byte, bc *Rwq_Blockchain) {
	var buff bytes.Buffer
	var payload rwq_version
	// 获取数据
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		log.Panic(err)
	}

	// 获取当前节点的最后height
	myBestHeight := bc.Rwq_GetBestHeight()
	// 数据中的最后height
	foreignerBestHeight := payload.Rwq_BestHeight

	// 节点将从消息中提取的 BestHeight 与自身进行比较
	// 当前的height比对方的小
	// 发送获取区块的的命令
	if myBestHeight < foreignerBestHeight {
		rwq_sendGetBlocks(payload.Rwq_AddrFrom)
	} else if myBestHeight > foreignerBestHeight {
		// 当前的height比对方的大
		// 通知对方节点，同步数据
		rwq_sendVersion(payload.Rwq_AddrFrom, bc)
	}

	// 如果不是已知节点，将节点添加到已知节点中
	if !rwq_nodeIsKnown(payload.Rwq_AddrFrom) {
		knownNodes = append(knownNodes, payload.Rwq_AddrFrom)
	}
}

// 是否是已知节点
func rwq_nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}

	return false
}

//================================================================

// 命令转字节数组
func rwq_commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

// 将字节数组转字符串命令
func rwq_bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return fmt.Sprintf("%s", command)
}

// 加密
func rwq_gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
