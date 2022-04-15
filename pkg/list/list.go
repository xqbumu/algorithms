package list

type Node struct {
	Val  int
	Next *Node
}

func New(vals ...int) *Node {
	if len(vals) == 0 {
		return nil
	}
	return &Node{
		Val:  vals[0],
		Next: New(vals[1:]...),
	}
}

func (n *Node) Extra() []int {
	return Extra(n)
}

func Extra(node *Node) []int {
	if node == nil {
		return []int{}
	}

	return append([]int{node.Val}, Extra(node.Next)...)
}
