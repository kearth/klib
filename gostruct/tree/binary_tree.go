package tree

import "fmt"

type TreeNode struct {
	Val   any
	Left  *TreeNode
	Right *TreeNode
}

func BuildTree(tl []any) *TreeNode {
	t := new(TreeNode)
	if len(tl) == 0 {
		return t
	}
	if tl[0] != nil {
		t.Val = tl[0]
		if 1 < len(tl) {
			t.Left = new(TreeNode)
			buildNode(t, tl, 0, 1)
		}
		if 2 < len(tl) {
			t.Right = new(TreeNode)
			buildNode(t, tl, 0, 2)
		}
	}
	return t
}

func buildNode(n *TreeNode, tl []any, level int, idx int) {
	if idx <= len(tl) && n != nil && tl[idx] != nil {
		n.Val = tl[idx]
		level = level + 1
		if level*2+1 < len(tl) {
			n.Left = new(TreeNode)
			buildNode(n.Left, tl, level, level*2+1)
		}
		if level*2+2 < len(tl) {
			n.Right = new(TreeNode)
			buildNode(n.Right, tl, level, level*2+2)
		}
	}
}

func (t *TreeNode) Print() {
	nodePrint(t)
}

func nodePrint(n *TreeNode) {
	fmt.Println(n.Val)
	if n.Left != nil {
		nodePrint(n.Left)
	}
	if n.Right != nil {
		nodePrint(n.Right)
	}
}

func Print(t *TreeNode) {
	nodePrint(t)
}

func MergeTrees(root1 *TreeNode, root2 *TreeNode) *TreeNode {
	root := new(TreeNode)
	mergeNode(root, root1, root2)
	return root
}

func mergeNode(ret *TreeNode, n1 *TreeNode, n2 *TreeNode) {
	if n1 != nil && n2 != nil {
		ret.Val = n1.Val.(int) + n2.Val.(int)
		mergeNode(ret.Left, n1.Left, n1.Left)
		mergeNode(ret.Right, n1.Right, n1.Right)
	} else if n1 != nil && n2 == nil {
		ret = n1
	} else if n1 == nil && n2 != nil {
		ret = n2
	} else {
		ret = nil
	}
}
