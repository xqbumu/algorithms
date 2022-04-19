package question

import (
	"log"
	"testing"
)

// https://leetcode-cn.com/problems/valid-sudoku/
func isValidSudoku(board [][]byte) bool {
	const size = 9
	nums := make([]byte, size)
	var i, j, sj, si, c, x, y int8

	for i = 0; i < size; i++ {
		if !checkBlock(board[i]...) {
			return false
		}
		for j = 0; j < size; j++ {
			nums[j] = board[j][i]
		}
		if !checkBlock(nums...) {
			return false
		}
		c, x, y = 0, i/3, i%3
		for si = 0; si < 3; si++ {
			for sj = 0; sj < 3; sj++ {
				nums[c] = board[x*3+si][y*3+sj]
				c++
			}
		}
		if !checkBlock(nums...) {
			return false
		}
	}
	return true
}

func checkBlock(nums ...byte) bool {
	var bits int16
	for i := 0; i < len(nums); i++ {
		if nums[i] == '.' {
			continue
		}
		if bits&(1<<(nums[i]-'1')) != 0 {
			printBlock(nums...)
			return false
		}
		bits |= 1 << (nums[i] - '1')
	}

	return true
}

func printBlock(nums ...byte) {
	for i := 0; i < 3; i++ {
		log.Printf("%s\n", nums[i*3:(i+1)*3])
	}
}

func Test_isValidSudoku(t *testing.T) {
	type args struct {
		board [][]byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{``, args{[][]byte{
			{'5', '3', '.', '.', '7', '.', '.', '.', '.'},
			{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
			{'.', '9', '8', '.', '.', '.', '.', '6', '.'},
			{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
			{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
			{'7', '.', '.', '.', '2', '.', '.', '.', '6'},
			{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
			{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
			{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
		}}, true},
		{``, args{[][]byte{
			{'8', '3', '.', '.', '7', '.', '.', '.', '.'},
			{'6', '.', '.', '1', '9', '5', '.', '.', '.'},
			{'.', '9', '8', '.', '.', '.', '.', '6', '.'},
			{'8', '.', '.', '.', '6', '.', '.', '.', '3'},
			{'4', '.', '.', '8', '.', '3', '.', '.', '1'},
			{'7', '.', '.', '.', '2', '.', '.', '.', '6'},
			{'.', '6', '.', '.', '.', '.', '2', '8', '.'},
			{'.', '.', '.', '4', '1', '9', '.', '.', '5'},
			{'.', '.', '.', '.', '8', '.', '.', '7', '9'},
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidSudoku(tt.args.board); got != tt.want {
				t.Errorf("isValidSudoku() = %v, want %v", got, tt.want)
			}
		})
	}
}
