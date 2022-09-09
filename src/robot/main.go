package main

import (
	"context"
	"fmt"
	_ "github.com/gorilla/websocket"
	"log"
	"os"
	"projectName/pkg/client"
	"projectName/pkg/dto"
	"projectName/pkg/solia"
	"projectName/pkg/token"
	"strings"
)

func main() {
	//机器人appid 和 token
	token := token.BotToken(0000, "0000")
	var api client.OpenAPI
	api.Token = token
	api.SetupClient()
	ctx := context.Background()
	ws, err := api.WS(ctx, nil, "") //websocket
	if err != nil {
		log.Fatalln("websocket错误， err = ", err)
		os.Exit(1)
	}
	start(token, ws, &api)

}

//开始
func start(token *token.Token, ws *dto.WebsocketAP, api *client.OpenAPI) {
	session := dto.Session{
		URL:     "wss://api.sgroup.qq.com/websocket/",
		Token:   *token,
		Intent:  dto.IntentGuildAtMessage, //固定了监听@机器人的消息
		LastSeq: 0,
		Shards: dto.ShardConfig{
			ShardID:    0,
			ShardCount: ws.Shards,
		},
	}
	var client client.Client
	wsClient := client.New(session)
	if err := wsClient.Connect(); err != nil {
		fmt.Println(err)
		return
	}
	var err error
	// 如果 session id 不为空，则执行的是 resume 操作，如果为空，则执行的是 identify 操作
	if session.ID != "" {
		err = wsClient.Resume()
	} else {
		// 初次鉴权
		err = wsClient.Identify()
	}
	if err != nil {
		log.Println(fmt.Sprintf("[ws/session] Identify/Resume err %+v", err))
		return
	}
	ctx := context.Background()

	var so solia.Solia
	var handel dto.ATMessageEventHandler = func(event1 *dto.WSPayload, data *dto.Message) error {
		if strings.Index(data.Content, "hello") > -1 || strings.Index(data.Content, "你好") > -1 { // 如果@机器人并输入 hello or 你好 则回复 你好。
			api.PostMessage(ctx, data.ChannelID, data.Author.ID, &dto.MessageToCreate{MsgID: data.ID, Content: solia.Hello})
		}
		if strings.Index(data.Content, solia.Solitaire) > -1 { // 如果@机器人并输入 成语接龙 则开始游戏。
			str, err := so.ReadStart(data.Author.ID)
			var msg string
			if err != nil {
				msg = err.Error()
			} else {
				msg = solia.Start + "『" + str + "』"
			}
			api.PostMessage(ctx, data.ChannelID, data.Author.ID, &dto.MessageToCreate{MsgID: data.ID, Content: msg})
		} else if so.UserId != "" && so.UserId == data.Author.ID && strings.Index(data.Content, solia.Cancel) == -1 { //开始接龙
			str, err := so.ReadStr(data.Content)
			var msg string
			if err != nil {
				msg = err.Error()
			} else {
				msg = solia.Action + "『" + str + "』"
			}
			api.PostMessage(ctx, data.ChannelID, data.Author.ID, &dto.MessageToCreate{MsgID: data.ID, Content: msg})
		} else if so.UserId != "" && so.UserId == data.Author.ID && strings.Index(data.Content, solia.Cancel) > -1 {
			so.ReNew()
			api.PostMessage(ctx, data.ChannelID, data.Author.ID, &dto.MessageToCreate{MsgID: data.ID, Content: "游戏结束"})
		}
		return nil
	}
	dto.DefaultHandlers.ATMessage = handel
	wsClient.Listening()
}
