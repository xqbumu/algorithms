package main

import (
	"algorithms/pkg/list"
	"reflect"
	"testing"
)

type ListNode = list.Node

// https://leetcode-cn.com/problems/reverse-nodes-in-k-group/
func reverseKGroup(head *ListNode, k int) *ListNode {
	var res, pre, sub *ListNode
	var dRes, dHead, dNext *ListNode

	do2 := func() *ListNode {
		dHead, dRes = pre, pre
		for i := 0; i < k-1; i++ {
			dNext = dRes
			dRes = dHead.Next
			dHead.Next = dRes.Next
			dRes.Next = dNext
		}
		return dRes
	}

	var mod, pos int
	for head != nil {
		mod = pos % k
		if mod == 0 {
			pre = head
		}
		if mod == k-1 {
			if res == nil {
				res = do2()
			} else {
				sub.Next = do2()
			}
			sub = pre
			head = pre.Next
			pos++
			continue
		}
		head = head.Next
		pos++
	}

	return res
}

func Test_reverseKGroup(t *testing.T) {
	type args struct {
		head *ListNode
		k    int
	}
	tests := []struct {
		name string
		args args
		want *ListNode
	}{
		{``, args{list.New(1, 2, 3), 2}, list.New(2, 1, 3)},
		{``, args{list.New(1, 2, 3, 4), 2}, list.New(2, 1, 4, 3)},
		{``, args{list.New(1, 2, 3, 4, 5), 1}, list.New(1, 2, 3, 4, 5)},
		{``, args{list.New(1, 2, 3, 4, 5), 2}, list.New(2, 1, 4, 3, 5)},
		{``, args{list.New(1, 2, 3, 4, 5, 6, 7), 2}, list.New(2, 1, 4, 3, 6, 5, 7)},
		{``, args{list.New(1, 2, 3, 4, 5), 3}, list.New(3, 2, 1, 4, 5)},
		{``, args{list.New(1, 2, 3, 4, 5), 4}, list.New(4, 3, 2, 1, 5)},
		{``, args{list.New(1), 1}, list.New(1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reverseKGroup(tt.args.head, tt.args.k)
			if !reflect.DeepEqual(got.Extra(), tt.want.Extra()) {
				t.Errorf("reverseKGroup() = %v, want %v", got.Extra(), tt.want.Extra())
			}
		})
	}
}
