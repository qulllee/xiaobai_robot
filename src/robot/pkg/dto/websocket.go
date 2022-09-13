package dto

import "projectName/pkg/token"

// Intent 类型
type Intent int

const (
	IntentGuilds         Intent = 1 << iota
	IntentGuildAtMessage Intent = 1 << 30 // 只接收@消息事件
)

// EventType 事件类型
type EventType string

// OPCode websocket op 码
type OPCode int

type Session struct {
	ID      string
	URL     string
	Token   token.Token
	Intent  Intent
	LastSeq uint32
	Shards  ShardConfig
}

// ShardConfig 连接的 shard 配置，ShardID 从 0 开始，ShardCount 最小为 1
type ShardConfig struct {
	ShardID    uint32
	ShardCount uint32
}

// WSUser 当前连接的用户信息
type WSUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Bot      bool   `json:"bot"`
}

// WSPayload websocket 消息结构
type WSPayload struct {
	WSPayloadBase
	Data       interface{} `json:"d,omitempty"`
	RawMessage []byte      `json:"-"` // 原始的 message 数据
}

// WebsocketAP wss 接入点信息
type WebsocketAP struct {
	URL               string            `json:"url"`
	Shards            uint32            `json:"shards"`
	SessionStartLimit SessionStartLimit `json:"session_start_limit"`
}

// SessionStartLimit 链接频控信息
type SessionStartLimit struct {
	Total          uint32 `json:"total"`
	Remaining      uint32 `json:"remaining"`
	ResetAfter     uint32 `json:"reset_after"`
	MaxConcurrency uint32 `json:"max_concurrency"`
}

// WSResumeData 重连数据
type WSResumeData struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       uint32 `json:"seq"`
}

// WSPayloadBase 基础消息结构，排除了 data
type WSPayloadBase struct {
	OPCode OPCode    `json:"op"`
	Seq    uint32    `json:"s,omitempty"`
	Type   EventType `json:"t,omitempty"`
}

// WS OPCode
const (
	WSDispatchEvent OPCode = iota
	WSHeartbeat
	WSIdentity
	_ // Presence Update
	_ // Voice State Update
	_
	WSResume
	WSReconnect
	_ // Request Guild Members
	WSInvalidSession
	WSHello
	WSHeartbeatAck
	HTTPCallbackAck
)

// opMeans op 对应的含义字符串标识
var opMeans = map[OPCode]string{
	WSDispatchEvent:  "Event",
	WSHeartbeat:      "Heartbeat",
	WSIdentity:       "Identity",
	WSResume:         "Resume",
	WSReconnect:      "Reconnect",
	WSInvalidSession: "InvalidSession",
	WSHello:          "Hello",
	WSHeartbeatAck:   "HeartbeatAck",
}

// OPMeans 返回 op 含义
func OPMeans(op OPCode) string {
	means, ok := opMeans[op]
	if !ok {
		means = "unknown"
	}
	return means
}

// WSHelloData hello 返回
type WSHelloData struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

// WSIdentityData 鉴权数据
type WSIdentityData struct {
	Token      string   `json:"token"`
	Intents    Intent   `json:"intents"`
	Shard      []uint32 `json:"shard"` // array of two integers (shard_id, num_shards)
	Properties struct {
		Os      string `json:"$os,omitempty"`
		Browser string `json:"$browser,omitempty"`
		Device  string `json:"$device,omitempty"`
	} `json:"properties,omitempty"`
}
