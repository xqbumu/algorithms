package tree

type Node struct {
	Val   int
	Left  *Node
	Right *Node
}

func IndexLevel(i int, fix bool) int {
	if fix {
		i += 1
	}
	var c = 0
	for i > 0 {
		i >>= 1
		c++
	}
	return c
}

func New(vals ...interface{}) *Node {
	if len(vals) == 0 {
		return nil
	}
	val, ok := vals[0].(int)
	if !ok {
		return nil
	}
	root := &Node{
		Val: val,
	}

	for pos := 1; pos < len(vals); pos++ {
		v, ok := vals[pos].(int)
		if !ok {
			continue
		}

		n, _ := FindNodeByPos(root, nil, pos, IndexLevel(pos, true)-1, true)
		n.Val = v
	}

	return root
}

func FindNodeByPos(root, parent *Node, pos int, lv int, ac bool) (*Node, *Node) {
	if root == nil || lv == 0 || pos == 0 {
		return root, nil
	}

	switch ((pos + 1) >> (lv - 1)) & 1 {
	case 0:
		if root.Left == nil && ac {
			root.Left = &Node{}
		}
		return FindNodeByPos(root.Left, root, pos, lv-1, ac)
	case 1:
		if root.Right == nil && ac {
			root.Right = &Node{}
		}
		return FindNodeByPos(root.Right, root, pos, lv-1, ac)
	}

	return nil, nil
}
