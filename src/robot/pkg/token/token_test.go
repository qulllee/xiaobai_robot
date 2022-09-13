package token

import (
	"reflect"
	"testing"
)

func TestBotToken(t *testing.T) {
	type args struct {
		appID       uint64
		accessToken string
	}
	tests := []struct {
		name string
		args args
		want *Token
	}{
		{
			name: "test1-生成token",
			args: args{
				appID:       1232,
				accessToken: "1232",
			},
			want: &Token{
				AppID:       1232,
				AccessToken: "1232",
				Type:        "Bot",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BotToken(tt.args.appID, tt.args.accessToken); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BotToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
