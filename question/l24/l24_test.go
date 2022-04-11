package l24

import (
	"reflect"
	"testing"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func NewListNode(vals ...int) *ListNode {
	if len(vals) == 0 {
		return nil
	}
	return &ListNode{
		Val:  vals[0],
		Next: NewListNode(vals[1:]...),
	}
}

func ExtraListNode(node *ListNode) []int {
	if node == nil {
		return []int{}
	}

	return append([]int{node.Val}, ExtraListNode(node.Next)...)
}

func swapPairs(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}
	a, b := head, head.Next
	a.Next, b.Next = swapPairs(b.Next), a

	return b
}

func Test_swapPairs(t *testing.T) {
	type args struct {
		head *ListNode
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{``, args{NewListNode(1, 2, 3, 4)}, NewListNode(2, 1, 4, 3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := swapPairs(tt.args.head)
			if !reflect.DeepEqual(ExtraListNode(tt.want), ExtraListNode(got)) {
				t.Errorf("swapPairs() = %v, want %v", ExtraListNode(got), ExtraListNode(tt.want))
			}
		})
	}
}
