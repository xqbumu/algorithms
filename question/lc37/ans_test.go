package question

import (
	"log"
	"reflect"
	"testing"
)

// https://leetcode-cn.com/problems/sudoku-solver/
func solveSudoku(board [][]byte) {
	const size = 9
	var i, j, sj, si, c, x, y int8
	rows := map[byte][]byte{
		'x': make([]byte, 9),
		'y': make([]byte, 9),
		'z': make([]byte, 9),
	}

	for i = 0; i < size; i++ {
		rows['x'] = append([]byte{}, board[0]...)
		for j = 0; j < size; j++ {
			rows['y'][j] = board[j][i]
		}
		c, x, y = 0, i/3, i%3
		for si = 0; si < 3; si++ {
			for sj = 0; sj < 3; sj++ {
				rows['z'][c] = board[x*3+si][y*3+sj]
				c++
			}
		}
		nums := checkRows(rows)
		printRows(rows)
		log.Printf("-- %s\n", nums)
	}
}

func checkRows(rows map[byte][]byte) []byte {
	var bits int16
	for _, row := range rows {
		for i := 0; i < len(row); i++ {
			if row[i] == '.' {
				continue
			}
			if bits&(1<<(row[i]-'0')) != 0 {
				continue
			}
			bits |= 1 << (row[i] - '0')
		}
	}

	res := []byte{}
	for i := 0; i < 10; i++ {
		if bits>>i&1 != 0 {
			continue
		}
		res = append(res, byte('0'+i))
	}

	return res
}

func printRows(rows map[byte][]byte) {
	for _, row := range rows {
		for i := 0; i < 3; i++ {
			log.Printf("%s\n", row[i*3:(i+1)*3])
		}
	}
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
			if reflect.DeepEqual(tt.args.board, tt.want) {
				t.Errorf("solveSudoku() = %v, want %v", tt.args.board, tt.want)
			}
		})
	}
}
