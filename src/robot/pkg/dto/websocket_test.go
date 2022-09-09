package dto

import "testing"

func TestOPMeans(t *testing.T) {
	type args struct {
		op OPCode
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				op: 0,
			},
			want: "Event",
		},
		{
			name: "test1",
			args: args{
				op: -1,
			},
			want: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OPMeans(tt.args.op); got != tt.want {
				t.Errorf("OPMeans() = %v, want %v", got, tt.want)
			}
		})
	}
}
