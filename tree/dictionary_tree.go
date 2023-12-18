package tree

import (
	"dict-go/utils"
	"fmt"
)

type DictionaryTree interface {
	Iterate(ch chan *WordNode)
	Insert(word string, meaning string) *WordNode
	Remove(word string) *WordNode
	Search(word string) *WordNode
	createNode(word string, meaning string) *WordNode
}

type WordNode struct {
	left    *WordNode
	Word    string
	Meaning string
	right   *WordNode
}

func (node *WordNode) Iterate(ch chan *WordNode) {
	if node == nil {
		return
	}

	node.left.Iterate(ch)
	ch <- node
	node.right.Iterate(ch)
}

func (node *WordNode) SetData(data *WordNode) {
	node.Word = data.Word
	node.Meaning = data.Meaning
}

func (node *WordNode) Insert(word string, meaning string) *WordNode {
	if node == nil {
		return node.createNode(word, meaning)
	}

	newDataHash := utils.ToDJBHash(word)
	curDataHash := utils.ToDJBHash(node.Word)
	switch {
	case curDataHash > newDataHash:
		node.left = node.left.Insert(word, meaning)
	case curDataHash < newDataHash:
		node.right = node.right.Insert(word, meaning)
	}

	return node
}

func (node *WordNode) Remove(word string) *WordNode {
	if node == nil {
		return node
	}

	targetDataHash := utils.ToDJBHash(word)
	curDataHash := utils.ToDJBHash(node.Word)
	switch {
	case curDataHash > targetDataHash:
		node.left = node.left.Remove(word)
	case curDataHash < targetDataHash:
		node.right = node.right.Remove(word)
	default:
		switch {
		case node.left == nil:
			right := node.right
			if right != nil {
				node.right = nil
			}

			return right
		case node.right == nil:
			left := node.left
			if left != nil {
				node.left = nil
			}

			return left
		default:
			tmp := node.right
			for tmp != nil && tmp.left != nil {
				tmp = tmp.left
			}

			node.SetData(tmp)
			node.right = node.right.Remove(word)
		}
	}

	return node
}

func (node *WordNode) Search(word string) *WordNode {
	if node == nil {
		return node
	}

	targetDataHash := utils.ToDJBHash(word)
	curDataHash := utils.ToDJBHash(node.Word)

	if curDataHash == targetDataHash {
		return node
	}

	if curDataHash > targetDataHash {
		return node.left.Search(word)
	}

	return node.right.Search(word)
}

func (node *WordNode) createNode(word string, meaning string) *WordNode {
	return &WordNode{
		Word:    word,
		Meaning: meaning,
	}
}

func (node *WordNode) ModifyMeaning(newMeaning string) {
	node.Meaning = newMeaning
}

func (node *WordNode) String() string {
	return fmt.Sprintf("%s;%s", node.Word, node.Meaning)
}
