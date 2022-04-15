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

	var pi, qi int = -1, -1
	for pos := 0; pos < 100000; pos++ {
		n, _ := tree.FindNodeByPos(root, nil, pos, tree.IndexLevel(pos, true)-1, false)
		if n == nil {
			continue
		}
		if n.Val == p.Val && pi == -1 {
			pi = pos + 1
			// fmt.Println(pos, n.Val)
		}
		if n.Val == q.Val && qi == -1 {
			qi = pos + 1
			// fmt.Println(pos, n.Val)
		}
		if pi != -1 && qi != -1 {
			break
		}
	}
	// fmt.Println("pos:", pi, qi)

	min := tree.IndexLevel(pi, false)
	if qs := tree.IndexLevel(qi, false); qs < min {
		min = qs
		pi >>= (tree.IndexLevel(pi, false) - min)
	} else {
		qi >>= (tree.IndexLevel(qi, false) - min)
	}

	var pos int
	for i := min; i >= 0; i-- {
		if (qi>>i)&1 == (pi>>i)&1 {
			pos <<= 1
			pos += (qi >> i & 1)
		} else {
			break
		}
	}
	// fmt.Println("result:", min, pos)

	n, _ := tree.FindNodeByPos(root, nil, pos-1, tree.IndexLevel(pos-1, true)-1, false)

	return n
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
		{``, args{tree.New(data.List236[0]...), tree.New(17), tree.New(2)}, tree.New(2)},
		{``, args{tree.New(data.List236[0]...), tree.New(17), tree.New(3)}, tree.New(1)},
		{``, args{tree.New(data.List236[0]...), tree.New(17), tree.New(4)}, tree.New(4)},
		{``, args{tree.New(data.List236[0]...), tree.New(17), tree.New(20)}, tree.New(2)},
		{``, args{tree.New(3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4), tree.New(5), tree.New(4)}, tree.New(5)},
		{``, args{tree.New(3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4), tree.New(5), tree.New(1)}, tree.New(3)},
		{``, args{tree.New(3, 5, 1, 6, 2, 0, 8, nil, nil, 7, 4), tree.New(6), tree.New(4)}, tree.New(5)},
		{``, args{tree.New(1, 2), tree.New(1), tree.New(2)}, tree.New(1)},
		{``, args{tree.New(1, 2, 3, nil, 4), tree.New(4), tree.New(1)}, tree.New(1)},
		{``, args{tree.New(1, 2, 3, nil, 4), tree.New(4), tree.New(3)}, tree.New(1)},
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
