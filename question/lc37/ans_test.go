package question

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const dim = 9

// https://leetcode-cn.com/problems/sudoku-solver/
func solveSudoku(board [][]byte) {
	exec(board)
}

func exec(board [][]byte) bool {
	for x, y, opts, done := checkDone(board); !done; {
		for i := 0; i < len(opts); i++ {
			board[x][y] = opts[i]
			if exec(board) {
				return true
			}
			board[x][y] = '.'
		}
		return false
	}
	return true
}

func checkDone(board [][]byte) (int8, int8, []byte, bool) {
	var i, j, x, y int8
	var opts []byte
	done := true
	for i = 0; i < dim; i++ {
		for j = 0; j < dim; j++ {
			if board[i][j] != '.' {
				continue
			}
			rows := relateRows(board, i, j)
			nums := remainNums(rows)
			if done == true {
				done = false
				x, y, opts = i, j, nums
			} else if len(nums) < len(opts) {
				x, y, opts = i, j, nums
			}

		}
	}
	return x, y, opts, done
}

func relateRows(board [][]byte, x, y int8) [][]byte {
	var sy, sx, cnt, xx, yy int8
	rows := [][]byte{
		make([]byte, 9),
		make([]byte, 9),
		make([]byte, 9),
	}
	rows[0] = append([]byte{}, board[x]...)
	for t := 0; t < dim; t++ {
		rows[1][t] = board[t][y]
	}
	cnt, xx, yy = 0, x/3, y/3
	for sx = 0; sx < 3; sx++ {
		for sy = 0; sy < 3; sy++ {
			rows[2][cnt] = board[xx*3+sx][yy*3+sy]
			cnt++
		}
	}

	return rows
}

func remainNums(rows [][]byte) []byte {
	var bits uint16
	for _, row := range rows {
		for i := 0; i < len(row); i++ {
			if row[i] == '.' {
				continue
			}
			if bits&(1<<(row[i]-'0')) != 0 {
				continue
			}
			bits |= (1 << (row[i] - '0'))
		}
	}

	res := []byte{}
	for i := 1; i < 10; i++ {
		if bits>>i&1 != 0 {
			continue
		}
		res = append(res, byte('0'+i))
	}

	return res
}

func printBoard(board [][]byte) {
	out := bytes.NewBuffer(nil)
	for i := 0; i < dim; i++ {
		fmt.Fprintf(out, "%s\n", board[i])
	}
	log.Printf("\n%s", out)
}

func printRows(rows [][]byte) {
	out := bytes.NewBuffer(nil)
	for _, row := range rows {
		for i := 0; i < 3; i++ {
			fmt.Fprintf(out, "%s\n", row[i*3:(i+1)*3])
		}
	}
	log.Printf("\n%s", out)
}

func Test_solveSudoku(t *testing.T) {
	type args struct {
		board [][]byte
	}
	tests := []struct {
		name string
		args args
		want [][]byte
	}{
		{``, args{[][]byte{
			{'5', '3', '.', '6', '7', '8', '9', '1', '2'},
			{'6', '7', '2', '1', '.', '5', '3', '4', '8'},
			{'1', '9', '8', '3', '4', '2', '5', '6', '7'},
			{'8', '5', '9', '7', '6', '1', '4', '2', '3'},
			{'4', '2', '6', '8', '5', '3', '7', '9', '1'},
			{'7', '1', '3', '9', '2', '4', '8', '5', '6'},
			{'9', '6', '1', '5', '3', '7', '2', '8', '4'},
			{'2', '8', '7', '4', '1', '9', '6', '3', '5'},
			{'3', '4', '5', '2', '8', '6', '1', '7', '9'},
		}}, [][]byte{
			{'5', '3', '4', '6', '7', '8', '9', '1', '2'},
			{'6', '7', '2', '1', '9', '5', '3', '4', '8'},
			{'1', '9', '8', '3', '4', '2', '5', '6', '7'},
			{'8', '5', '9', '7', '6', '1', '4', '2', '3'},
			{'4', '2', '6', '8', '5', '3', '7', '9', '1'},
			{'7', '1', '3', '9', '2', '4', '8', '5', '6'},
			{'9', '6', '1', '5', '3', '7', '2', '8', '4'},
			{'2', '8', '7', '4', '1', '9', '6', '3', '5'},
			{'3', '4', '5', '2', '8', '6', '1', '7', '9'},
		}},
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
		}}, [][]byte{
			{'5', '3', '4', '6', '7', '8', '9', '1', '2'},
			{'6', '7', '2', '1', '9', '5', '3', '4', '8'},
			{'1', '9', '8', '3', '4', '2', '5', '6', '7'},
			{'8', '5', '9', '7', '6', '1', '4', '2', '3'},
			{'4', '2', '6', '8', '5', '3', '7', '9', '1'},
			{'7', '1', '3', '9', '2', '4', '8', '5', '6'},
			{'9', '6', '1', '5', '3', '7', '2', '8', '4'},
			{'2', '8', '7', '4', '1', '9', '6', '3', '5'},
			{'3', '4', '5', '2', '8', '6', '1', '7', '9'},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			solveSudoku(tt.args.board)
			assert.Equal(t, tt.want, tt.args.board)
		})
	}
}
