package goadb

import (
	"reflect"
	"testing"
)

func Test_splitCmdAgrs(t *testing.T) {
	type args struct {
		cmds string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{"1", args{cmds: "ab bc cd"}, []string{"ab", "bc", "cd"}},
		{"2", args{cmds: "  ab   bc   cd "}, []string{"ab", "bc", "cd"}},
		{"3", args{cmds: "ab b\\ c cd"}, []string{"ab", "b c", "cd"}},
		{"4", args{cmds: " ab  b\\\\c   cd"}, []string{"ab", "b\\\\c", "cd"}},
		{"5", args{cmds: " ab  b\\\\c   c\\ 你好d"}, []string{"ab", "b\\\\c", "c 你好d"}},
		{"6", args{cmds: ""}, []string{}},
		{"7", args{cmds: "      "}, []string{}},
		{"8", args{cmds: "   \\   "}, []string{" "}},
		{"9", args{cmds: "ab\\cd"}, []string{"ab\\cd"}},
		{"10", args{cmds: "\\ab cd"}, []string{"\\ab", "cd"}},
		{"11", args{cmds: "\\ab cd\\"}, []string{"\\ab", "cd\\"}},
		{"12", args{cmds: "\\ ab cd"}, []string{" ab", "cd"}},
		{"13", args{cmds: "ab cd\\ "}, []string{"ab", "cd "}},
		{"14", args{cmds: "ab cd\\    "}, []string{"ab", "cd "}},
		{"15", args{cmds: "ab cd\\ \\\\   "}, []string{"ab", "cd \\ "}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitCmdAgrs(tt.args.cmds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitCmdAgrs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_safeArg(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"1", args{arg: "ab"}, "ab"},
		{"2", args{arg: "  "}, "\\ \\ "},
		{"3", args{arg: "a b"}, "a\\ b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := safeArg(tt.args.arg); got != tt.want {
				t.Errorf("safeArg() = %v, want %v", got, tt.want)
			}
		})
	}
}
