package cmd

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantArgs Args
	}{
		{"MoveToStart", "0", Args{Mv: '0'}},
		{"10 x MoveDown", "10j", Args{Has: true, Num: 10, Mv: 'j'}},
		{"MoveRight", "l", Args{Num: 1, Mv: 'l'}},
		{"MoveToLastLine", "G", Args{Num: 1, Mv: 'G'}},
		{"change word", "cw", Args{Num: 1, Op: 'c', SubNum: 1, Mv: 'w'}},
		{"delete line", "dd", Args{Num: 1, Op: 'd', SubNum: 1, Mv: 'd'}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArgs := Parse([]rune(tt.input))
			if gotArgs != tt.wantArgs {
				t.Errorf(
					"Parse([]rune(\"%s\") = %v; wanted %v",
					tt.input, gotArgs, tt.wantArgs,
				)
			}
		})
	}
}
