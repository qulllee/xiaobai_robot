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

func TestOpenAPI_WS(t *testing.T) {
	type fields struct {
		Token       *token.Token
		timeout     time.Duration
		Sandbox     bool
		debug       bool
		lastTraceID string
		restyClient *resty.Client
	}
	type args struct {
		ctx context.Context
		in1 map[string]string
		in2 string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dto.WebsocketAP
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				Token: &token.Token{
					AppID:       1232,
					AccessToken: "string",
					Type:        "bot",
				},
				timeout:     30,
				Sandbox:     false,
				debug:       false,
				lastTraceID: "string",
				restyClient: &resty.Client{},
			},
			args: args{
				ctx: context.Background(),
				in1: nil,
				in2: "",
			},
			want:    &dto.WebsocketAP{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OpenAPI{
				Token:       tt.fields.Token,
				timeout:     tt.fields.timeout,
				Sandbox:     tt.fields.Sandbox,
				debug:       tt.fields.debug,
				lastTraceID: tt.fields.lastTraceID,
				restyClient: tt.fields.restyClient,
			}
			got, err := o.WS(tt.args.ctx, tt.args.in1, tt.args.in2)
			if (err != nil) != tt.wantErr {
				t.Errorf("WS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WS() got = %v, want %v", got, tt.want)
			}
		})
	}
}

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
			name: "test1",
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
			name: "test2",
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
				msg:       &dto.MessageToCreate{MsgID: "08c9dff9c1b8fb8edec50110b2eb8d05389c0648eefdea9907", Content: solia.Hello},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test2",
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
				t.Errorf("PostMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
