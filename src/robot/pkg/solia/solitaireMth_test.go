package solia

import (
	"bufio"
	"testing"
)

func TestSolia_ReadStart(t *testing.T) {
	type fields struct {
		UserId string
		StrSet Set
		flag   bool
		rd     *bufio.Reader
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
			name: "test1",
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
			want:    "人山人海",
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
			got, err := s.ReadStart(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadStart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != "" && len(got) == len(tt.want) {
				t.Errorf("ReadStart() got = %v, want %v", got, tt.want)
			}
		})
	}
}
