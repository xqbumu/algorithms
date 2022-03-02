// 机器人坐标问题
// 问题描述
// 有一个机器人，给一串指令，L左转 R右转，F前进一步，B后退一步，问最后机器人的坐标，最开始，机器人位于 0 0，方向为正Y。 可以输入重复指令n ： 比如 R2(LF) 这个等于指令 RLFLF。 问最后机器人的坐标是多少？
package question

import (
	"reflect"
	"strconv"
	"testing"
)

type position struct {
	x  int
	y  int
	dx int
	dy int
}

func newPosition() *position {
	return &position{
		dy: 1,
	}
}

func (p *position) left() {
	if p.dx == 1 {
		p.dx, p.dy = 0, 1
	} else if p.dx == -1 {
		p.dx, p.dy = 0, -1
	} else if p.dy == 1 {
		p.dx, p.dy = -1, 0
	} else if p.dy == -1 {
		p.dx, p.dy = 1, 0
	}
}

func (p *position) right() {
	if p.dx == 1 {
		p.dx, p.dy = 0, -1
	} else if p.dx == -1 {
		p.dx, p.dy = 0, 1
	} else if p.dy == 1 {
		p.dx, p.dy = 1, 0
	} else if p.dy == -1 {
		p.dx, p.dy = -1, 0
	}
}

func (p *position) forward() {
	p.x, p.y = p.x+p.dx, p.y+p.dy
}

func (p *position) back() {
	p.x, p.y = p.x-p.dx, p.y-p.dy
}

// move 移动搁置
func (p *position) move(cmd string) (np position) {
	mp := *p
	rPos := 0
	rNum := 1

	for i := 0; i < len(cmd); i++ {
		switch cmd[i] {
		case 'L': // 左转
			mp.left()
		case 'R': // 右转
			mp.right()
		case 'F': // 前进一步
			mp.forward()
		case 'B': // 后退一步
			mp.back()
		case '(':
			rPos = i - 1
			rNum -= 1
			continue
		case ')':
			if rNum > 0 {
				i = rPos
			}
			continue
		default:
			if n, err := strconv.Atoi(string(cmd[i])); err == nil && n > 0 {
				rNum = n
			}
		}
	}

	return mp
}

func Test_position_move(t *testing.T) {
	type args struct {
		cmd string
	}
	tests := []struct {
		name   string
		p      *position
		args   args
		wantNp position
	}{
		// {``, newPosition(), args{`R`}, position{0, 0, 1, 0}},
		// {``, newPosition(), args{`RF`}, position{1, 0, 1, 0}},
		{``, newPosition(), args{`R2(F)`}, position{2, 0, 1, 0}},
		{``, newPosition(), args{`R2(LF)`}, position{-1, 1, -1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNp := tt.p.move(tt.args.cmd)
			if !reflect.DeepEqual(gotNp, tt.wantNp) {
				t.Errorf("position.move() = %v, want %v", gotNp, tt.wantNp)
			}
		})
	}
}
