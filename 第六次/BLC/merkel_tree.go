package BLC

import "crypto/sha256"

type Rwq_MerkelTree struct {
	Rwq_RootNode *Rwq_MerkelNode
}

type Rwq_MerkelNode struct {
	Rwq_Left  *Rwq_MerkelNode
	Rwq_Right *Rwq_MerkelNode
	Rwq_Data  []byte
}

func Rwq_NewMerkelTree(data [][]byte) *Rwq_MerkelTree {
	var nodes []Rwq_MerkelNode

	// 如果交易数据不是双数，将最后一个交易复制添加到最后
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	// 生成所有的一级节点，存储到node中
	for _, dataum := range data {
		node := Rwq_NewMerkelNode(nil, nil, dataum)
		nodes = append(nodes, *node)
	}

	// 遍历生成顶层节点
	for i := 0;i<len(data)/2 ;i++{
		var newLevel []Rwq_MerkelNode
		for j:=0 ; j<len(nodes) ;j+=2  {
			node := Rwq_NewMerkelNode(&nodes[j],&nodes[j+1],nil)
			newLevel = append(newLevel,*node)
		}
		nodes = newLevel
	}

	//for ; len(nodes)==1 ;{
	//	var newLevel []Rwq_MerkelNode
	//	for j:=0 ; j<len(nodes) ;j+=2  {
	//		node := Rwq_NewMerkelNode(&nodes[j],&nodes[j+1],nil)
	//		newLevel = append(newLevel,*node)
	//	}
	//	nodes = newLevel
	//}
	mTree := Rwq_MerkelTree{&nodes[0]}
	return &mTree
}

// 新叶节点
func Rwq_NewMerkelNode(left, right *Rwq_MerkelNode, data []byte) *Rwq_MerkelNode {
	mNode := Rwq_MerkelNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Rwq_Data = hash[:]
	} else {
		prevHashes := append(left.Rwq_Data, right.Rwq_Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Rwq_Data = hash[:]
	}

	mNode.Rwq_Left = left
	mNode.Rwq_Right = right

	return &mNode
}
