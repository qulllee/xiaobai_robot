package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/tencent-connect/botgo/version"
	"github.com/tidwall/gjson"
	"log"
	"net"
	"net/http"
	"os"
	"projectName/pkg/dto"
	"projectName/pkg/token"
	"projectName/pkg/ver"
	"runtime"
	"sync"
	"time"
)

type closeErrorChan chan error

type messageChan chan *dto.WSPayload

// WSUser 当前连接的用户信息
type WSUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Bot      bool   `json:"bot"`
}

type Client struct {
	version         int
	conn            *websocket.Conn
	messageQueue    messageChan
	session         *dto.Session
	user            *WSUser
	closeChan       closeErrorChan
	heartBeatTicker *time.Ticker // 用于维持定时心跳
}

const DefaultQueueSize = 10000

// MaxIdleConns 默认指定空闲连接池大小
const MaxIdleConns = 3000

type OpenAPI struct {
	Token   *token.Token
	timeout time.Duration

	Sandbox     bool   // 请求沙箱环境
	debug       bool   // debug 模式，调试sdk时候使用
	lastTraceID string // lastTraceID id

	restyClient *resty.Client // resty client 复用
}

// TraceIDKey 机器人openapi返回的链路追踪ID
const TraceIDKey = "X-Tps-trace-ID"

func (o *OpenAPI) SetupClient() {
	o.restyClient = resty.New().
		SetTransport(createTransport(nil, MaxIdleConns)). // 自定义 transport
		SetDebug(o.debug).
		SetTimeout(o.timeout).
		SetAuthToken(o.Token.GetString()).
		SetAuthScheme(string(o.Token.Type)).
		SetHeader("User-Agent", version.String()).
		SetPreRequestHook(
			func(client *resty.Client, request *http.Request) error {
				// 执行请求前过滤器
				// 由于在 `OnBeforeRequest` 的时候，request 还没生成，所以 filter 不能使用，所以放到 `PreRequestHook`
				return DoReqFilterChains(request, nil)
			},
		).
		// 设置请求之后的钩子，打印日志，判断状态码
		OnAfterResponse(
			func(client *resty.Client, resp *resty.Response) error {
				fmt.Println(fmt.Sprintf("%v", respInfo(resp)))
				// 执行请求后过滤器
				if err := DoRespFilterChains(resp.Request.RawRequest, resp.RawResponse); err != nil {
					return err
				}
				traceID := resp.Header().Get(TraceIDKey)
				o.lastTraceID = traceID
				// 非成功含义的状态码，需要返回 error 供调用方识别
				if !IsSuccessStatus(resp.StatusCode()) {
					return errors.New(fmt.Sprintf("%d%v%s", resp.StatusCode(), resp.Body(), traceID))
				}
				return nil
			},
		)
}

// IsSuccessStatus 是否是成功的状态码
func IsSuccessStatus(code int) bool {
	if _, ok := successStatusSet[code]; ok {
		return true
	}
	return false
}

var successStatusSet = map[int]bool{
	http.StatusOK:        true,
	http.StatusNoContent: true,
}

// WS 获取带分片 WSS 接入点
func (o *OpenAPI) WS(ctx context.Context, _ map[string]string, _ string) (*dto.WebsocketAP, error) {
	resp, err := o.request(ctx).
		SetResult(dto.WebsocketAP{}).
		Get(ver.GetURL(ver.GatewayBotURI, o.Sandbox))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.WebsocketAP), nil
}

// request 每个请求，都需要创建一个 request
func (o *OpenAPI) request(ctx context.Context) *resty.Request {
	return o.restyClient.R().SetContext(ctx)
}

// PostMessage 发消息
func (o *OpenAPI) PostMessage(ctx context.Context, channelID string, userID string, msg *dto.MessageToCreate) (*dto.Message, error) {
	if userID != "" {
		msg.Content = fmt.Sprintf("<@%s>", userID) + msg.Content
	}
	resp, err := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetBody(msg).
		Post(ver.GetURL(ver.MessagesURI, o.Sandbox))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Message), nil
}

// PostSettingGuide 发送设置引导消息, atUserID为要at的用户
func (o *OpenAPI) PostSettingGuide(ctx context.Context,
	channelID string, atUserIDs []string) (*dto.Message, error) {
	var content string
	for _, userID := range atUserIDs {
		content += fmt.Sprintf("<@%s>", userID)
	}
	msg := &dto.SettingGuideToCreate{
		Content: content,
	}
	resp, err := o.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetBody(msg).
		Post(ver.GetURL(ver.SettingGuideURI, o.Sandbox))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*dto.Message), nil
}

// New 新建一个连接对象
func (c *Client) New(session dto.Session) Client {
	return Client{
		messageQueue:    make(messageChan, DefaultQueueSize),
		session:         &session,
		closeChan:       make(closeErrorChan, 10),
		heartBeatTicker: time.NewTicker(60 * time.Second), // 先给一个默认 ticker，在收到 hello 包之后，会 reset
	}
}

// Connect 连接到 websocket
func (c *Client) Connect() error {
	if c.session.URL == "" {
		return errors.New("ws ap url is invalid")
	}
	var err error
	c.conn, _, err = websocket.DefaultDialer.Dial(c.session.URL, nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v, connect err: %v", c.session, err))
		return err
	}
	log.Println(fmt.Sprintf("%v, url %s, connected", c.session, c.session.URL))

	return nil
}

func (c *Client) Listening() error {
	defer c.Close()
	go c.readMessageToQueue()
	go c.listenMessageAndHandle()

	// 接收 resume signal
	resumeSignal := make(chan os.Signal, 1)

	// handler message
	for {
		select {
		case <-resumeSignal: // 使用信号量控制连接立即重连
			log.Println(fmt.Sprintf("%v, received resumeSignal signal", c.session))
			return errors.New("need reconnect")
		case err := <-c.closeChan:
			// 关闭连接的错误码 https://bot.q.qq.com/wiki/develop/api/gateway/error/error.html
			fmt.Println(fmt.Sprintf("%v Listening stop. err is %v", c.session, err))
			// 不能够 identify 的错误
			//if websocket.IsCloseError(err, 4914, 4915) {
			//	err = errs.New(errs.CodeConnCloseCantIdentify, err.Error())
			//}
			//// 这里用 UnexpectedCloseError，如果有需要排除在外的 close error code，可以补充在第二个参数上
			//// 4009: session time out, 发了 reconnect 之后马上关闭连接时候的错误码，这个是允许 resumeSignal 的
			//if websocket.IsUnexpectedCloseError(err, 4009) {
			//	err = errs.New(errs.CodeConnCloseCantResume, err.Error())
			//}
			return err
		case <-c.heartBeatTicker.C:
			log.Println(fmt.Sprintf("%v listened heartBeat", c.session))
			heartBeatEvent := &dto.WSPayload{
				WSPayloadBase: dto.WSPayloadBase{
					OPCode: dto.WSHeartbeat,
				},
				Data: c.session.LastSeq,
			}
			// 不处理错误，Write 内部会处理，如果发生发包异常，会通知主协程退出
			_ = c.Write(heartBeatEvent)
		}
	}
}

// Write 往 ws 写入数据
func (c *Client) Write(message *dto.WSPayload) error {
	m, _ := json.Marshal(message)
	log.Println(fmt.Sprintf("%v write %s message, %v", c.session, dto.OPMeans(message.OPCode), string(m)))

	if err := c.conn.WriteMessage(websocket.TextMessage, m); err != nil {
		fmt.Println(fmt.Sprintf("%v WriteMessage failed, %v", c.session, err))
		c.closeChan <- err
		return err
	}
	return nil
}

// Close 关闭连接
func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		fmt.Println(fmt.Sprintf("%v, close conn err: %v", c.session, err))
	}
	c.heartBeatTicker.Stop()
}

func (c *Client) readMessageToQueue() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println(fmt.Sprintf("%v read message failed, %v, message %s", c.session, err, string(message)))
			close(c.messageQueue)
			c.closeChan <- err
			return
		}
		payload := &dto.WSPayload{}
		if err := json.Unmarshal(message, payload); err != nil {
			fmt.Println(fmt.Sprintf("%v json failed, %v", c.session, err))
			continue
		}
		payload.RawMessage = message
		fmt.Println(fmt.Sprintf("%v receive %s message, %s", c.session, dto.OPMeans(payload.OPCode), string(message)))
		if c.isHandleBuildIn(payload) {
			continue
		}
		c.messageQueue <- payload
	}
}

// isHandleBuildIn 内置的事件处理，处理那些不需要业务方处理的事件
// return true 的时候说明事件已经被处理了
func (c *Client) isHandleBuildIn(payload *dto.WSPayload) bool {
	switch payload.OPCode {
	case dto.WSHello: // 接收到 hello 后需要开始发心跳
		c.startHeartBeatTicker(payload.RawMessage)
	case dto.WSHeartbeatAck: // 心跳 ack 不需要业务处理
	case dto.WSReconnect: // 达到连接时长，需要重新连接，此时可以通过 resume 续传原连接上的事件
		c.closeChan <- errors.New("need reconnect")
	case dto.WSInvalidSession: // 无效的 sessionLog，需要重新鉴权
		c.closeChan <- errors.New("invalid session")
	default:
		return false
	}
	return true
}

// startHeartBeatTicker 启动定时心跳
func (c *Client) startHeartBeatTicker(message []byte) {
	helloData := &dto.WSHelloData{}
	if err := parseData(message, helloData); err != nil {
		fmt.Println(fmt.Sprintf("%v hello data parse failed, %v, message %v", c.session, err, message))
	}
	// 根据 hello 的回包，重新设置心跳的定时器时间
	c.heartBeatTicker.Reset(time.Duration(helloData.HeartbeatInterval) * time.Millisecond)
}
func parseData(message []byte, target interface{}) error {
	data := gjson.Get(string(message), "d")
	return json.Unmarshal([]byte(data.String()), target)
}

func (c *Client) listenMessageAndHandle() {
	defer func() {
		// panic，一般是由于业务自己实现的 handle 不完善导致
		// 打印日志后，关闭这个连接，进入重连流程
		if err := recover(); err != nil {
			PanicHandler(err, c.session)
			c.closeChan <- fmt.Errorf("panic: %v", err)
		}
	}()
	for payload := range c.messageQueue {
		c.saveSeq(payload.Seq)
		// ready 事件需要特殊处理
		if payload.Type == "READY" {
			continue
		}
		// 解析具体事件，并投递给业务注册的 handler
		if err := dto.ParseAndHandle(payload); err != nil {
			fmt.Println(fmt.Sprintf("%v parseAndHandle failed, %v", c.session, payload))
		}
	}
	log.Println(fmt.Sprintf("%v message queue is closed", c.session))
}

// ParseAndHandle 处理回调事件
//func ParseAndHandle(payload *dto.WSPayload) error {
//	// 指定类型的 handler
//	if h, ok := eventParseFuncMap[payload.OPCode][payload.Type]; ok {
//		return h(payload, payload.RawMessage)
//	}
//	// 透传handler，如果未注册具体类型的 handler，会统一投递到这个 handler
//	if DefaultHandlers.Plain != nil {
//		return DefaultHandlers.Plain(payload, payload.RawMessage)
//	}
//	return nil
//}

func (c *Client) saveSeq(seq uint32) {
	if seq > 0 {
		c.session.LastSeq = seq
	}
}

// PanicBufLen Panic 堆栈大小
var PanicBufLen = 1024

// PanicHandler 处理websocket场景的 panic ，打印堆栈
func PanicHandler(e interface{}, session *dto.Session) {
	buf := make([]byte, PanicBufLen)
	buf = buf[:runtime.Stack(buf, false)]
	fmt.Println(fmt.Sprintf("[PANIC]%v\n%v\n%s\n", session, e, buf))
}

func createTransport(localAddr net.Addr, idleConns int) *http.Transport {
	dialer := &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 60 * time.Second,
	}
	if localAddr != nil {
		dialer.LocalAddr = localAddr
	}
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          idleConns,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   idleConns,
		MaxConnsPerHost:       idleConns,
	}
}

// respInfo 用于输出日志的时候格式化数据
func respInfo(resp *resty.Response) string {
	bodyJSON, _ := json.Marshal(resp.Request.Body)
	return fmt.Sprintf(
		"[OPENAPI]%v %v, traceID:%v, status:%v, elapsed:%v req: %v, resp: %v",
		resp.Request.Method,
		resp.Request.URL,
		resp.Header().Get(TraceIDKey),
		resp.Status(),
		resp.Time(),
		string(bodyJSON),
		string(resp.Body()),
	)
}

// HTTPFilter 请求过滤器
type HTTPFilter func(req *http.Request, response *http.Response) error

var (
	filterLock         = sync.RWMutex{}
	reqFilterChainSet  = map[string]HTTPFilter{}
	reqFilterChains    []string
	respFilterChainSet = map[string]HTTPFilter{}
	respFilterChains   []string
)

// DoReqFilterChains 按照注册顺序执行请求过滤器
func DoReqFilterChains(req *http.Request, resp *http.Response) error {
	for _, name := range reqFilterChains {
		if _, ok := reqFilterChainSet[name]; !ok {
			continue
		}
		if err := reqFilterChainSet[name](req, resp); err != nil {
			return err
		}
	}
	return nil
}

// DoRespFilterChains 按照注册顺序执行返回过滤器
func DoRespFilterChains(req *http.Request, resp *http.Response) error {
	for _, name := range respFilterChains {
		if _, ok := respFilterChainSet[name]; !ok {
			continue
		}
		if err := respFilterChainSet[name](req, resp); err != nil {
			return err
		}
	}
	return nil
}

// Resume 重连
func (c *Client) Resume() error {
	payload := &dto.WSPayload{
		Data: &dto.WSResumeData{
			Token:     c.session.Token.GetString(),
			SessionID: c.session.ID,
			Seq:       c.session.LastSeq,
		},
	}
	payload.OPCode = dto.WSResume // 内嵌结构体字段，单独赋值
	return c.Write(payload)
}

// Identify 对一个连接进行鉴权，并声明监听的 shard 信息
func (c *Client) Identify() error {
	// 避免传错 intent
	if c.session.Intent == 0 {
		c.session.Intent = dto.IntentGuilds
	}
	payload := &dto.WSPayload{
		Data: &dto.WSIdentityData{
			Token:   c.session.Token.GetString(),
			Intents: c.session.Intent,
			Shard: []uint32{
				c.session.Shards.ShardID,
				c.session.Shards.ShardCount,
			},
		},
	}
	payload.OPCode = dto.WSIdentity
	return c.Write(payload)
}
