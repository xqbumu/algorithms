package question

import (
	"algorithms/pkg/tree"
	"algorithms/question/data"
	"testing"
)

func lowestCommonAncestor(root, p, q *tree.Node) *tree.Node {
	if root == nil {
		return nil
	}
	if root.Val == p.Val || root.Val == q.Val {
		return root
	}
	ln := lowestCommonAncestor(root.Left, p, q)
	rn := lowestCommonAncestor(root.Right, p, q)
	if ln != nil && rn != nil {
		return root
	}
	if ln == nil {
		return rn
	}

	return ln
}

func Test_lowestCommonAncestor(t *testing.T) {
	type args struct {
		root *tree.Node
		p    *tree.Node
		q    *tree.Node
	}
	tests := []struct {
		name string
		args args
		want *tree.Node
	}{
		{``, args{tree.New(3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4, nil, nil, nil, 9, nil, nil, nil, nil, 10), tree.New(5), tree.New(1)}, tree.New(3)},
		{``, args{tree.New(3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4), tree.New(5), tree.New(4)}, tree.New(5)},
		{``, args{tree.New(3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4), tree.New(5), tree.New(1)}, tree.New(3)},
		{``, args{tree.New(3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4), tree.New(6), tree.New(4)}, tree.New(5)},
		{``, args{tree.New(1, 2), tree.New(1), tree.New(2)}, tree.New(1)},
		{``, args{tree.New(data.List236[1]...), tree.New(9998), tree.New(9999)}, tree.New(155)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lowestCommonAncestor(tt.args.root, tt.args.p, tt.args.q)
			if got.Val != tt.want.Val {
				t.Errorf("lowestCommonAncestor() = %v, want %v", got.Val, tt.want.Val)
			}
		})
	}
}
