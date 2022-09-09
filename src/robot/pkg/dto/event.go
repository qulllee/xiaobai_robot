package dto

import (
	"encoding/json"
	"github.com/tidwall/gjson"
)

// 事件类型
const (
	EventGuildCreate           EventType = "GUILD_CREATE"
	EventGuildUpdate           EventType = "GUILD_UPDATE"
	EventGuildDelete           EventType = "GUILD_DELETE"
	EventChannelCreate         EventType = "CHANNEL_CREATE"
	EventChannelUpdate         EventType = "CHANNEL_UPDATE"
	EventChannelDelete         EventType = "CHANNEL_DELETE"
	EventGuildMemberAdd        EventType = "GUILD_MEMBER_ADD"
	EventGuildMemberUpdate     EventType = "GUILD_MEMBER_UPDATE"
	EventGuildMemberRemove     EventType = "GUILD_MEMBER_REMOVE"
	EventMessageCreate         EventType = "MESSAGE_CREATE"
	EventMessageReactionAdd    EventType = "MESSAGE_REACTION_ADD"
	EventMessageReactionRemove EventType = "MESSAGE_REACTION_REMOVE"
	EventAtMessageCreate       EventType = "AT_MESSAGE_CREATE"
	EventPublicMessageDelete   EventType = "PUBLIC_MESSAGE_DELETE"
	EventDirectMessageCreate   EventType = "DIRECT_MESSAGE_CREATE"
	EventDirectMessageDelete   EventType = "DIRECT_MESSAGE_DELETE"
	EventAudioStart            EventType = "AUDIO_START"
	EventAudioFinish           EventType = "AUDIO_FINISH"
	EventAudioOnMic            EventType = "AUDIO_ON_MIC"
	EventAudioOffMic           EventType = "AUDIO_OFF_MIC"
	EventMessageAuditPass      EventType = "MESSAGE_AUDIT_PASS"
	EventMessageAuditReject    EventType = "MESSAGE_AUDIT_REJECT"
	EventMessageDelete         EventType = "MESSAGE_DELETE"
	EventForumThreadCreate     EventType = "FORUM_THREAD_CREATE"
	EventForumThreadUpdate     EventType = "FORUM_THREAD_UPDATE"
	EventForumThreadDelete     EventType = "FORUM_THREAD_DELETE"
	EventForumPostCreate       EventType = "FORUM_POST_CREATE"
	EventForumPostDelete       EventType = "FORUM_POST_DELETE"
	EventForumReplyCreate      EventType = "FORUM_REPLY_CREATE"
	EventForumReplyDelete      EventType = "FORUM_REPLY_DELETE"
	EventForumAuditResult      EventType = "FORUM_PUBLISH_AUDIT_RESULT"
	EventInteractionCreate     EventType = "INTERACTION_CREATE"
)

var eventParseFuncMap = map[OPCode]map[EventType]eventParseFunc{
	WSDispatchEvent: {
		EventAtMessageCreate: atMessageHandler,
	},
}

type eventParseFunc func(event *WSPayload, message []byte) error

// ParseAndHandle 处理回调事件
func ParseAndHandle(payload *WSPayload) error {
	// 指定类型的 handler,目前只有写了机器人消息
	if h, ok := eventParseFuncMap[payload.OPCode][payload.Type]; ok {
		return h(payload, payload.RawMessage)
	}
	return nil
}

func atMessageHandler(payload *WSPayload, message []byte) error {
	data := &Message{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if payload != nil {
		return DefaultHandlers.ATMessage(payload, data)
	}
	return nil
}

// ParseData 解析数据
func ParseData(message []byte, target interface{}) error {
	data := gjson.Get(string(message), "d")
	return json.Unmarshal([]byte(data.String()), target)
}

//func ATMessageEventHandler(event *WSPayload, data *Message) error {
//	log.Println("12321321")
//	if strings.Index(data.Content,"hello")>-1 || strings.Index(data.Content,"你好")>-1{
//
//	}
//	return nil
//}
