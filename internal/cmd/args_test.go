package cmd

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     Args
		wantPair Pair
		wantOk   bool
	}{
		{"MoveToStart", Args{Mv: '0'}, Pair{Mv: Cmd{Kind: MoveToStart}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPair, gotOk := tt.args.Parse()
			if gotPair != tt.wantPair || gotOk != tt.wantOk {
				t.Errorf(
					"%v.Parse() = %v, %v; wanted %v, %v",
					tt.args, gotPair, gotOk, tt.wantPair, tt.wantOk,
				)
			}
		})
	}
}
