package solia

import (
	"testing"
)

func TestMpSolia_ReadStart(t *testing.T) {
	type fields struct {
		rd []string
		Mp map[string]*Solia
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1-1号选手开始成语接龙",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼"},
				Mp: map[string]*Solia{},
			},
			args: args{
				userID: "001",
			},
			want:    "人山人海",
			wantErr: false,
		},
		{
			name: "test2-1号选手重复开始成语接龙(1号选手游戏中)",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼"},
				Mp: map[string]*Solia{
					"001": &Solia{},
				},
			},
			args: args{
				userID: "001",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test3-2号选手开始成语接龙(1号选手游戏中)",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼"},
				Mp: map[string]*Solia{
					"001": &Solia{},
				},
			},
			args: args{
				userID: "002",
			},
			want:    "人山人海",
			wantErr: false,
		},
		{
			name: "test4-2号选手重复开始成语接龙(1、2号选手游戏中)",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼"},
				Mp: map[string]*Solia{
					"001": &Solia{},
					"002": &Solia{},
				},
			},
			args: args{
				userID: "002",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MpSolia{
				rd: tt.fields.rd,
				Mp: tt.fields.Mp,
			}
			got, err := ms.ReadStart(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadStart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadStart() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMpSolia_readLineNum(t *testing.T) {
	type fields struct {
		rd []string
		Mp map[string]*Solia
	}
	type args struct {
		lineNum int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1-获取第2行的成语(词库长度之类)",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼", "天外有天", "人外有人"},
				Mp: map[string]*Solia{},
			},
			args: args{
				lineNum: 2,
			},
			want:    "海市蜃楼",
			wantErr: false,
		},
		{
			name: "test2-获取第8行的成语(超出词库长度)",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼", "天外有天", "人外有人"},
				Mp: map[string]*Solia{},
			},
			args: args{
				lineNum: 8,
			},
			want:    "人山人海",
			wantErr: false,
		},
		{
			name: "test3-获取第-1行的成语(行号小于等于0)",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼", "天外有天", "人外有人"},
				Mp: map[string]*Solia{},
			},
			args: args{
				lineNum: -1,
			},
			want:    "人山人海",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MpSolia{
				rd: tt.fields.rd,
				Mp: tt.fields.Mp,
			}
			got, err := ms.readLineNum(tt.args.lineNum)
			if (err != nil) != tt.wantErr {
				t.Errorf("readLineNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readLineNum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMpSolia_ReadStr(t *testing.T) {
	type fields struct {
		rd []string
		Mp map[string]*Solia
	}
	type args struct {
		content string
		userId  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1-1号游戏者正常接龙(输入正确成语)",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼", "天外有天", "人外有人", "楼外青山"},
				Mp: map[string]*Solia{
					"001": &Solia{tryNum: 3, nowStr: "人山人海", StrSet: Set{m: map[interface{}]struct{}{"人山人海": struct{}{}}}},
				},
			},
			args: args{
				content: "海市蜃楼",
				userId:  "001",
			},
			want:    "楼外青山",
			wantErr: false,
		},
		{
			name: "test2-1号游戏者正常接龙(输入正确成语,机器人接不出)",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼", "天外有天", "人外有人"},
				Mp: map[string]*Solia{
					"001": &Solia{tryNum: 3, nowStr: "人山人海", StrSet: Set{m: map[interface{}]struct{}{"人山人海": struct{}{}}}},
				},
			},
			args: args{
				content: "海市蜃楼",
				userId:  "001",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test3-1号游戏者正常接龙(输入错误成语,还剩三次机会)",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼", "天外有天", "人外有人", "楼外青山"},
				Mp: map[string]*Solia{
					"001": &Solia{tryNum: 3, nowStr: "人山人海", StrSet: Set{m: map[interface{}]struct{}{"人山人海": struct{}{}}}},
				},
			},
			args: args{
				content: "哈哈哈哈",
				userId:  "001",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test4-1号游戏者正常接龙(输入错误成语,还剩二次机会)",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼", "天外有天", "人外有人", "楼外青山"},
				Mp: map[string]*Solia{
					"001": &Solia{tryNum: 2, nowStr: "人山人海", StrSet: Set{m: map[interface{}]struct{}{"人山人海": struct{}{}}}},
				},
			},
			args: args{
				content: "哈哈哈哈",
				userId:  "001",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test5-1号游戏者正常接龙(输入错误成语,还剩一次机会)",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼", "天外有天", "人外有人", "楼外青山"},
				Mp: map[string]*Solia{
					"001": &Solia{tryNum: 1, nowStr: "人山人海", StrSet: Set{m: map[interface{}]struct{}{"人山人海": struct{}{}}}},
				},
			},
			args: args{
				content: "哈哈哈哈",
				userId:  "001",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test6-2号游戏者正常接龙(1号游戏者正在游戏)",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼", "天外有天", "人外有人", "楼外青山"},
				Mp: map[string]*Solia{
					"001": &Solia{tryNum: 1, nowStr: "人山人海", StrSet: Set{m: map[interface{}]struct{}{"人山人海": struct{}{}}}},
					"002": &Solia{tryNum: 1, nowStr: "人山人海", StrSet: Set{m: map[interface{}]struct{}{"人山人海": struct{}{}}}},
				},
			},
			args: args{
				content: "海市蜃楼",
				userId:  "002",
			},
			want:    "楼外青山",
			wantErr: false,
		},
		{
			name: "test7-2号游戏者错误接龙(1号游戏者正在游戏)",
			fields: fields{
				rd: []string{"人山人海", "海市蜃楼", "天外有天", "人外有人", "楼外青山"},
				Mp: map[string]*Solia{
					"001": &Solia{tryNum: 1, nowStr: "人山人海", StrSet: Set{m: map[interface{}]struct{}{"人山人海": struct{}{}}}},
					"002": &Solia{tryNum: 1, nowStr: "人山人海", StrSet: Set{m: map[interface{}]struct{}{"人山人海": struct{}{}}}},
				},
			},
			args: args{
				content: "海里有余",
				userId:  "002",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MpSolia{
				rd: tt.fields.rd,
				Mp: tt.fields.Mp,
			}
			got, err := ms.ReadStr(tt.args.content, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadStr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
