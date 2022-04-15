package l24

import (
	"algorithms/pkg/list"
	"reflect"
	"testing"
)

func swapPairs(head *list.Node) *list.Node {
	if head == nil || head.Next == nil {
		return head
	}
	a, b := head, head.Next
	a.Next, b.Next = swapPairs(b.Next), a

	return b
}

func Test_swapPairs(t *testing.T) {
	type args struct {
		head *list.Node
	}
	tests := []struct {
		name string
		args args
		want *list.Node
	}{
		{``, args{list.New(1, 2, 3, 4)}, list.New(2, 1, 4, 3)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := swapPairs(tt.args.head)
			if !reflect.DeepEqual(list.Extra(tt.want), list.Extra(got)) {
				t.Errorf("swapPairs() = %v, want %v", list.Extra(got), list.Extra(tt.want))
			}
		})
	}
}
