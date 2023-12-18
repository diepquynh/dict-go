package main

import (
	"bufio"
	"dict-go/tree"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var FILENAME string = "dict.txt"
var HELP_STRING string = `
Tu dien Golang
*. Thoat khoi chuong trinh
1. Nhap du lieu tu dien tu file
2. Tim kiem 1 tu
3. Thay doi nghia 1 tu
4. Xoa 1 tu trong tu dien
5. Xuat cac tu da co trong tu dien ra man hinh
6. Xuat cac tu da co trong tu dien ra file
Moi nhap lua chon: `

func check(e error) int {
	if e != nil {
		return -1
	}

	return 0
}

func preImport(root tree.DictionaryTree, b []byte) tree.DictionaryTree {
	dataStr := string(b)
	for _, data := range strings.Split(dataStr, "\n") {
		if len(data) == 0 {
			continue
		}

		dataNode := strings.Split(data, ";")
		root = root.Insert(dataNode[0], dataNode[1])
	}

	return root
}

func importFromFile(root tree.DictionaryTree, filename string) tree.DictionaryTree {
	data, err := os.ReadFile(filename)
	if check(err) != 0 {
		return nil
	}

	return preImport(root, data)
}

func exportToFile(root tree.DictionaryTree, filename string) (ret int) {
	f, err := os.Create(filename)
	ret = check(err)
	if ret != 0 {
		return
	}

	defer f.Close()

	wordChannel := make(chan *tree.WordNode)
	go func(channel chan *tree.WordNode) {
		root.Iterate(channel)
		close(channel)
	}(wordChannel)

	for {
		wordNode, ok := <-wordChannel
		if !ok {
			return
		}

		f.WriteString(wordNode.String() + "\n")
	}
}

func iterateAndPrint(root tree.DictionaryTree) {
	wordChannel := make(chan *tree.WordNode)
	go func(channel chan *tree.WordNode) {
		root.Iterate(channel)
		close(channel)
	}(wordChannel)

	for {
		wordNode, ok := <-wordChannel
		if !ok {
			return
		}

		printNode(wordNode)
	}
}

func printNode(node *tree.WordNode) {
	fmt.Println("Tu:", node.Word)
	fmt.Println("Nghia:", node.Meaning)
}

func waitForEnter(reader *bufio.Reader) {
	fmt.Println("Nhan enter de tiep tuc")
	reader.ReadBytes('\n')
}

func main() {
	var root tree.DictionaryTree
	var node *tree.WordNode
	root = node
	if root = importFromFile(root, FILENAME); root != nil {
		fmt.Printf("Du lieu co san ton tai. Da tien hanh nhap du lieu tu dien co san\n\n")
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(HELP_STRING)
		line, _, _ := reader.ReadLine()
		choice, _ := strconv.Atoi(string(line))

		switch choice {
		case 1:
			if root = importFromFile(root, FILENAME); root != nil {
				fmt.Println("Nhap du lieu thanh cong")
			}

			waitForEnter(reader)
		case 2:
			fmt.Print("Nhap tu can tim kiem: ")
			bytes, _, _ := reader.ReadLine()
			word := string(bytes)

			node := root.Search(word)
			if node == nil {
				fmt.Printf("Tu %s khong ton tai\n", word)
				waitForEnter(reader)
				break
			}

			printNode(node)
			waitForEnter(reader)
		case 3:
			fmt.Print("Nhap tu can tim kiem: ")
			bytes, _, _ := reader.ReadLine()
			word := string(bytes)

			node := root.Search(word)
			if node == nil {
				fmt.Printf("Tu %s khong ton tai\n", word)
				waitForEnter(reader)
				break
			}

			printNode(node)
			fmt.Print("Nhap nghia moi: (enter de bo qua)")
			bytes, _, _ = reader.ReadLine()
			newMeaning := string(bytes)
			if len(newMeaning) > 0 {
				node.ModifyMeaning(newMeaning)
			}
			waitForEnter(reader)
		case 4:
			fmt.Print("Nhap tu can xoa: ")
			bytes, _, _ := reader.ReadLine()
			word := string(bytes)
			root.Remove(word)
			fmt.Println("Xoa tu thanh cong")
			waitForEnter(reader)
		case 5:
			iterateAndPrint(root)
			waitForEnter(reader)
		case 6:
			if exportToFile(root, FILENAME) == 0 {
				fmt.Println("Xuat du lieu ra file thanh cong!")
			}
			waitForEnter(reader)
		default:
			fmt.Println("Dang thoat khoi chuong trinh...")
			return
		}
	}
}
