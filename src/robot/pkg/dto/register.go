package dto

var DefaultHandlers struct {
	ATMessage ATMessageEventHandler
}

// ATMessageEventHandler at 机器人消息事件 handler
type ATMessageEventHandler func(event *WSPayload, data *Message) error
