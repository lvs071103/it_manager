package main

import "fmt"

// 方法
type TreeNode struct {
	Left, Right *TreeNode
	Value       int
}

// 工厂函数 值引用
func CreateNode(value int) *TreeNode {
	// 返回局部变量的地址
	return &TreeNode{Value: value}
}

// node为接收者
func (node TreeNode) print() {
	fmt.Println(node.Value)
}

// 地址引用
// nil 指针也可以调用方法
func (node *TreeNode) setValue(value int) {
	if node == nil {
		fmt.Println("Setting value to nil node. Ignored")
		return
	}
	node.Value = value
}

// 遍历树
func (node *TreeNode) traverse() {
	if node == nil {
		return
	}
	node.Left.traverse()
	node.print()
	node.Right.traverse()
}

func main() {
	var root TreeNode
	fmt.Println(root)

	root = TreeNode{Value: 3}
	root.Left = &TreeNode{}
	root.Right = &TreeNode{nil, nil, 5}
	root.Right.Left = new(TreeNode)
	root.Left.Right = CreateNode(2)
	root.traverse()

	nodes := []TreeNode{
		{Value: 3},
		{},
		{nil, &root, 6},
	}
	fmt.Println(nodes)
	root.print()
	root.Right.Left.setValue(4)
	root.Right.Left.print()
	pRoot := &root
	pRoot.print()
	pRoot.setValue(100)
	pRoot.print()

	var aRoot *TreeNode
	aRoot.setValue(200)
	aRoot = &root
	aRoot.setValue(300)
	aRoot.print()
	// root.traverse()

}
