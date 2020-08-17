package pkg

import (
	"errors"
	"fmt"
	"github.com/oldthreefeng/stress/utils"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	HttpOk         = 200 // 请求成功
	RequestTimeout = 506 // 请求超时
	RequestErr     = 509 // 请求错误
	ParseError     = 510 // 解析错误

	FormTypeHttp      = "http"
	FormTypeWebSocket = "webSocket"
	DefaultMethod     = "GET"
	DefaultVerifyCode = "statusCode"
	DefaultTimeOut    = 30 * time.Second
)

var (
	// 校验函数
	verifyMapHttp      = make(map[string]VerifyHttp)
	verifyMapHttpMutex sync.RWMutex

	verifyMapWebSocket      = make(map[string]VerifyWebSocket)
	verifyMapWebSocketMutex sync.RWMutex
)

// RegisterVerifyHttp is 注册http校验函数
func RegisterVerifyHttp(verify string, verifyFunc VerifyHttp) {
	verifyMapHttpMutex.Lock()
	defer verifyMapHttpMutex.Unlock()

	key := fmt.Sprintf("%s.%s", FormTypeHttp, verify)
	verifyMapHttp[key] = verifyFunc
}

// RegisterVerifyWebSocket is 注册webSocket校验函数
func RegisterVerifyWebSocket(verify string, verifyFunc VerifyWebSocket) {
	verifyMapWebSocketMutex.Lock()
	defer verifyMapWebSocketMutex.Unlock()

	key := fmt.Sprintf("%s.%s", FormTypeWebSocket, verify)
	verifyMapWebSocket[key] = verifyFunc
}

// Verify is interface for 验证器
type Verify interface {
	GetCode() int    // 有一个方法，返回code为200为成功
	GetResult() bool // 返回是否成功
}

// VerifyHttp is verify for http 验证方法
type VerifyHttp func(request *Request, response *http.Response) (code int, isSucceed bool)

// VerifyWebSocket is verify for websocket
type VerifyWebSocket func(request *Request, seq string, msg []byte) (code int, isSucceed bool)

// Request is a form for http request 请求结果
type Request struct {
	Url     string              // Url
	Form    string              // http/webSocket/tcp
	Method  string              // 方法 GET/POST/PUT
	Headers utils.ConcurrentMap // Headers
	Body    string              // body
	Verify  string              // 验证的方法
	Timeout time.Duration       // 请求超时时间
	Debug   bool                // 是否开启Debug模式

	// 连接以后初始化事件
	// 循环事件 切片 时间 动作
}

// GetBody return io.Reader from request
func (r *Request) GetBody() (body io.Reader) {
	body = strings.NewReader(r.Body)

	return
}

func (r *Request) getVerifyKey() (key string) {
	key = fmt.Sprintf("%s.%s", r.Form, r.Verify)

	return
}

// GetVerifyHttp is 获取数据校验方法
func (r *Request) GetVerifyHttp() VerifyHttp {
	verify, ok := verifyMapHttp[r.getVerifyKey()]
	if !ok {
		panic("GetVerifyHttp 验证方法不存在:" + r.Verify)
	}

	return verify
}

// GetVerifyWebSocket is 获取数据校验方法
func (r *Request) GetVerifyWebSocket() VerifyWebSocket {
	verify, ok := verifyMapWebSocket[r.getVerifyKey()]
	if !ok {
		panic("GetVerifyWebSocket 验证方法不存在:" + r.Verify)
	}

	return verify
}

// NewRequest is get Request form
// url 压测的url
// verify 验证方法 在server/verify中 http 支持:statusCode、json webSocket支持:json
// timeout 请求超时时间
// debug 是否开启debug
// path curl文件路径 http接口压测，自定义参数设置
func NewRequest(url string, verify string, timeout time.Duration, debug bool, reqHeaders []string, reqBody string) (request *Request, err error) {

	var (
		method  = DefaultMethod
		headers = utils.New(Concurrency)
		body    string
	)

	if reqBody != "" {
		method = "POST"
		body = reqBody
		headers.Set("Content-Type","application/x-www-form-urlencoded; charset=utf-8")
	}

	for _, v := range reqHeaders {
		getHeaderValue(v, headers)
	}

	form := ""
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		form = FormTypeHttp
	} else if strings.HasPrefix(url, "ws://") || strings.HasPrefix(url, "wss://") {
		form = FormTypeWebSocket
	} else {
		form = FormTypeHttp
		url = fmt.Sprintf("http://%s", url)
	}

	if form == "" {
		err = errors.New(fmt.Sprintf("url:%s 不合法,必须是完整http、webSocket连接", url))

		return
	}

	var (
		ok bool
	)

	switch form {
	case FormTypeHttp:
		// verify
		key := fmt.Sprintf("%s.%s", form, verify)
		_, ok = verifyMapHttp[key]
		if !ok {
			err = errors.New("验证器不存在:" + key)

			return
		}
	case FormTypeWebSocket:
		// verify
		if verify == DefaultVerifyCode {
			verify = "json"
		}

		key := fmt.Sprintf("%s.%s", form, verify)
		_, ok = verifyMapWebSocket[key]
		if !ok {
			err = errors.New("验证器不存在:" + key)

			return
		}

	}

	if timeout == 0 {
		timeout = DefaultTimeOut
	}

	request = &Request{
		Url:     url,
		Form:    form,
		Method:  strings.ToUpper(method),
		Headers: headers,
		Body:    body,
		Verify:  verify,
		Timeout: timeout,
		Debug:   debug,
	}

	return

}

// NewDefaultRequest is a default request
func NewDefaultRequest() *Request {
	return &Request{
		Url:     "http://www.baidu.com",
		Form:    FormTypeHttp,
		Method:  DefaultMethod,
		Verify:  VerifyStr,
		Timeout: DefaultTimeOut,
		Debug:   Debug,
		Body:    "",
	}
}

func getHeaderValue(v string, headers utils.ConcurrentMap) {
	index := strings.Index(v, ":")
	if index < 0 {
		return
	}

	vIndex := index + 1
	if len(v) >= vIndex {
		value := strings.TrimPrefix(v[vIndex:], " ")

		if val, ok := headers.Get(v[:index]); ok {
			headers.Set(v[:index], fmt.Sprintf("%s; %s", val, value))
		} else {
			headers.Set(v[:index], value)
		}

	}
}

// Print is 打印参数
func (r *Request) Print() {
	if r == nil {

		return
	}

	result := fmt.Sprintf("request:\n form:%s \n url:%s \n method:%s \n headers:%v \n", r.Form, r.Url, r.Method, r.Headers.Items())
	result = fmt.Sprintf("%s data:%v \n", result, r.Body)
	result = fmt.Sprintf("%s verify:%s \n timeout:%s \n debug:%v \n", result, r.Verify, r.Timeout, r.Debug)
	fmt.Println(result)

	return
}

// GetDebug is use Debug mode
func (r *Request) GetDebug() bool {

	return r.Debug
}

// IsParameterLegal is verify param is legal
func (r *Request) IsParameterLegal() (err error) {

	r.Form = "http"
	// statusCode json
	r.Verify = "json"

	key := fmt.Sprintf("%s.%s", r.Form, r.Verify)
	_, ok := verifyMapHttp[key]
	if !ok {

		return errors.New("验证器不存在:" + key)
	}

	return
}

// RequestResults is  请求结果
type RequestResults struct {
	Id        string // 消息Id
	ChanId    int // 消息Id
	Time      uint64 // 请求时间 纳秒
	IsSucceed bool   // 是否请求成功
	ErrCode   int    // 错误码
}

// SetId is set chanId & id to request results
func (r *RequestResults) SetId(chanId int, number uint64) {
	id := fmt.Sprintf("%d_%d", chanId, number)

	r.Id = id
	r.ChanId = chanId
}
