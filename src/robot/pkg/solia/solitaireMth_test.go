package solia

import (
	"testing"
)

func TestSolia_ReadStart(t *testing.T) {
	type fields struct {
		UserId string
		StrSet Set
		flag   bool
		rd     []string
		tryNum int
		nowStr string
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
			name: "test1(正常开始)", //正常开始，成语接龙
			fields: fields{
				UserId: "",
				StrSet: Set{
					m: make(map[interface{}]struct{}),
				},
				flag:   false,
				rd:     nil,
				tryNum: 3,
				nowStr: "",
			},
			args: args{
				userID: "1313",
			},
			want:    "人山人海",
			wantErr: false,
		},
		{
			name: "test2", //已经开始成语接龙，其他人再次调用
			fields: fields{
				UserId: "13131",
				StrSet: Set{
					m: make(map[interface{}]struct{}),
				},
				flag:   false,
				rd:     nil,
				tryNum: 3,
				nowStr: "",
			},
			args: args{
				userID: "1313",
			},
			want:    "人山人海",
			wantErr: true,
		},
		{
			name: "test3(重复开始)", //已经开始成语接龙，本人再次调用
			fields: fields{
				UserId: "13131",
				StrSet: Set{
					m: make(map[interface{}]struct{}),
				},
				flag:   false,
				rd:     nil,
				tryNum: 3,
				nowStr: "",
			},
			args: args{
				userID: "13131",
			},
			want:    "人山人海",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Solia{
				UserId: tt.fields.UserId,
				StrSet: tt.fields.StrSet,
				flag:   tt.fields.flag,
				rd:     tt.fields.rd,
				tryNum: tt.fields.tryNum,
				nowStr: tt.fields.nowStr,
			}
			got, err := s.ReadStart(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadStart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != "" && len(got) != len(tt.want) {
				t.Errorf("ReadStart() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSolia_readLineNum(t *testing.T) {
	type fields struct {
		UserId string
		StrSet Set
		flag   bool
		rd     []string
		tryNum int
		nowStr string
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
			name: "test1", //获取1000行后，第一个四字成语
			fields: fields{
				UserId: "13131",
				StrSet: Set{
					m: make(map[interface{}]struct{}),
				},
				flag:   false,
				rd:     nil,
				tryNum: 3,
				nowStr: "",
			},
			args: args{
				lineNum: 1000,
			},
			want:    "不足为训",
			wantErr: false,
		},
		{
			name: "test2", //获取-1行后，第一个四字成语
			fields: fields{
				UserId: "13131",
				StrSet: Set{
					m: make(map[interface{}]struct{}),
				},
				flag:   false,
				rd:     nil,
				tryNum: 3,
				nowStr: "",
			},
			args: args{
				lineNum: -1,
			},
			want:    "阿鼻地狱",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Solia{
				UserId: tt.fields.UserId,
				StrSet: tt.fields.StrSet,
				flag:   tt.fields.flag,
				rd:     tt.fields.rd,
				tryNum: tt.fields.tryNum,
				nowStr: tt.fields.nowStr,
			}
			got, err := s.readLineNum(tt.args.lineNum)
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

func TestSolia_ReadStr(t *testing.T) {
	type fields struct {
		UserId string
		StrSet Set
		flag   bool
		rd     []string
		tryNum int
		nowStr string
	}
	type args struct {
		content string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1", //接不出下一个成语
			fields: fields{
				UserId: "13131",
				StrSet: Set{
					m: make(map[interface{}]struct{}),
				},
				flag:   false,
				rd:     nil,
				tryNum: 3,
				nowStr: "声威大震",
			},
			args: args{
				content: "<@21321321> 震耳欲聋",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test2", //输入非成语
			fields: fields{
				UserId: "13131",
				StrSet: Set{
					m: make(map[interface{}]struct{}),
				},
				flag:   false,
				rd:     nil,
				tryNum: 3,
				nowStr: "声威大震",
			},
			args: args{
				content: "<@21321321> 乱写的",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test3", //正常接出下一个成语
			fields: fields{
				UserId: "13131",
				StrSet: Set{
					m: make(map[interface{}]struct{}),
				},
				flag:   false,
				rd:     nil,
				tryNum: 3,
				nowStr: "声威大震",
			},
			args: args{
				content: "<@21321321> 震天动地",
			},
			want:    "地主之谊",
			wantErr: false,
		},
		{
			name: "test4(首字和上一个成语不一样)", //首字和上一个成语不一样
			fields: fields{
				UserId: "13131",
				StrSet: Set{
					m: make(map[interface{}]struct{}),
				},
				flag:   false,
				rd:     nil,
				tryNum: 3,
				nowStr: "声威大震",
			},
			args: args{
				content: "<@21321321> 天崩地裂",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Solia{
				UserId: tt.fields.UserId,
				StrSet: tt.fields.StrSet,
				flag:   tt.fields.flag,
				rd:     tt.fields.rd,
				tryNum: tt.fields.tryNum,
				nowStr: tt.fields.nowStr,
			}
			got, err := s.ReadStr(tt.args.content)
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
