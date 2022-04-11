package question

import (
	"log"
	"testing"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func NewTreeNode(vals ...interface{}) *TreeNode {
	if len(vals) == 0 {
		return nil
	}
	val, ok := vals[0].(int)
	if !ok {
		return nil
	}
	root := &TreeNode{
		Val: val,
	}

	for pos := 1; pos < len(vals); pos++ {
		val, ok := vals[pos].(int)
		if !ok {
			continue
		}
		AppendValue(root, pos, 0, val)
	}

	return root
}

func AppendValue(node *TreeNode, pos int, lv int, val int) {
	switch pos {
	case 0:
		node.Left = &TreeNode{
			Val: val,
		}
	case 1:
		node.Right = &TreeNode{
			Val: val,
		}
	default:
		switch pos / 2 {
		case 0:
			AppendValue(node.Left, pos-(1<<lv+1), lv+1, val)
		case 1:
			AppendValue(node.Right, pos-(1<<lv+1), lv+1, val)
		}
	}
}

func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	return nil
}

func Test_lowestCommonAncestor(t *testing.T) {
	type args struct {
		head *TreeNode
	}
	tests := []struct {
		name string
		args args
		want *TreeNode
	}{
		{``, args{NewTreeNode(3, 5, 1, 6, 2, 0, 8)}, NewTreeNode(3, 5, 1, 6)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Println(tt)
			// got := lowestCommonAncestor(tt.args.head)
			// if !reflect.DeepEqual(ExtraTreeNode(tt.want), ExtraTreeNode(got)) {
			// 	t.Errorf("swapPairs() = %v, want %v", ExtraTreeNode(got), ExtraTreeNode(tt.want))
			// }
		})
	}
}
