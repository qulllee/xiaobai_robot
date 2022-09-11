package client

import (
	"context"
	"github.com/go-resty/resty/v2"
	"projectName/pkg/dto"
	"projectName/pkg/solia"
	"projectName/pkg/token"
	"reflect"
	"testing"
	"time"
)

func TestOpenAPI_PostMessage(t *testing.T) {
	type fields struct {
		Token       *token.Token
		timeout     time.Duration
		Sandbox     bool
		debug       bool
		lastTraceID string
		restyClient *resty.Client
	}
	type args struct {
		ctx       context.Context
		channelID string
		userID    string
		msg       *dto.MessageToCreate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.Message
		wantErr bool
	}{
		{
			name: "test1-回复用户打招呼信息(错误用户消息)",
			fields: fields{
				Token: &token.Token{
					AppID:       123,
					AccessToken: "213",
					Type:        "Bot",
				},
			},
			args: args{
				ctx:       context.Background(),
				channelID: "123",
				userID:    "123",
				msg:       &dto.MessageToCreate{MsgID: "123", Content: solia.Hello},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test2-回复用户打招呼信息",
			fields: fields{
				Token: &token.Token{
					AppID:       102019593,
					AccessToken: "qlABeIFRQz1WUXyyBW7gQwMcnaj27Y6x",
					Type:        "Bot",
				},
			},
			args: args{
				ctx:       context.Background(),
				channelID: "10712498",
				userID:    "6931576746694468327",
				msg:       &dto.MessageToCreate{Content: solia.Hello},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test3-发送下测试消息",
			fields: fields{
				Token: &token.Token{
					AppID:       102019593,
					AccessToken: "qlABeIFRQz1WUXyyBW7gQwMcnaj27Y6x",
					Type:        "Bot",
				},
			},
			args: args{
				ctx:       context.Background(),
				channelID: "10712498",
				userID:    "6931576746694468327",
				msg:       &dto.MessageToCreate{Content: "测试一下"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var o OpenAPI
			o.Token = tt.fields.Token
			o.SetupClient()
			o.WS(tt.args.ctx, nil, "")
			got, err := o.PostMessage(tt.args.ctx, tt.args.channelID, tt.args.userID, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostMessage() error = %s, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
