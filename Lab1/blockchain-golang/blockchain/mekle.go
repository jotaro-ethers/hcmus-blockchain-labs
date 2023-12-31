package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}

	node.Left = left
	node.Right = right

	return &node
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}

	for _, dat := range data {
		node := NewMerkleNode(nil, nil, dat)
		nodes = append(nodes, *node)
	}

	for i := 0; i < len(data)/2; i++ {
		var level []MerkleNode

		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			level = append(level, *node)
		}

		nodes = level
	}
	tree := MerkleTree{&nodes[0]}

	return &tree
}


func printMerkleNode(node *MerkleNode, level int, isLeft bool, isRoot bool) {
	if node == nil {
		return
	}

	indent := strings.Repeat("  ", level)

	fmt.Printf("%s", indent)
	if isLeft {
		fmt.Print("├─ ")
	} else {
		fmt.Print("└─ ")
	}
	fmt.Printf("%s\n", hex.EncodeToString(node.Data))

	if node.Left != nil || node.Right != nil {
		printMerkleNode(node.Left, level+1, true, false)
		printMerkleNode(node.Right, level+1, false, false)
	}
}

func PrintMerkleTree(block *Block) {
	transactions := make([][]byte, len(block.Transactions))
	for i, tx := range block.Transactions {
		transactions[i] = tx.Data
	}

	tree := NewMerkleTree(transactions)
	fmt.Println("Merkle Tree:")
	printMerkleNode(tree.RootNode, 0, true, true)
}




