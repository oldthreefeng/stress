package pkg

import (
	"github.com/oldthreefeng/stress/utils"
	"math/rand"
	"sync"
	"time"
)

// Http is http go link
func Http(chanId int, ch chan<- *RequestResults, totalNumber uint64, wg *sync.WaitGroup, request *Request) {

	defer func() {
		wg.Done()
	}()

	// fmt.Printf("启动协程 编号:%05d \n", chanId)
	for i := uint64(0); i < totalNumber; i++ {
		// 如果curl路径存在， 则处理curl。
		list := getRequestList(request, Path)

		requestResults := sendList(list)
		for _, v := range requestResults {
			v.SetId(chanId, i*uint64(len(list)))
			ch <- v
		}
	}

	return
}

// sendList 多个接口分步压测
func sendList(requestList []*Request) (requestResults []*RequestResults) {
	var (
		requestTime uint64
	)
	for _, request := range requestList {
		succeed, code, u := send(request)
		result := &RequestResults{
			Time:      requestTime + u,
			ErrCode:   code,
			IsSucceed: succeed,
		}
		requestResults = append(requestResults, result)
		if succeed == false {
			break
		}
	}

	return
}

// send 发送一次请求
func send(request *Request) (bool, int, uint64) {
	var (
		// startTime = time.Now()
		isSucceed = false
		errCode   = HttpOk
	)

	newRequest := getRequest(request)
	// newRequest := request

	resp, requestTime, err := HttpRequest(newRequest.Method, newRequest.Url, newRequest.GetBody(), newRequest.Headers, newRequest.Timeout)
	// requestTime := uint64(utils.DiffNano(startTime))
	if err != nil {
		errCode = RequestErr // 请求错误
	} else {
		// 验证请求是否成功
		errCode, isSucceed = newRequest.GetVerifyHttp()(newRequest, resp)
	}
	return isSucceed, errCode, requestTime
}

// ReqListWeigh  接口加权压测
type ReqListWeigh struct {
	list       []Req
	weighCount uint32 // 总权重
}

// Req is weights for request
type Req struct {
	req     *Request // 请求信息
	weights uint32   // 权重，数字越大访问频率越高
}

func (r *ReqListWeigh) setWeighCount() {
	r.weighCount = 0
	for _, value := range r.list {
		r.weighCount = r.weighCount + value.weights
	}
}

var (
	clientWeigh *ReqListWeigh
	r           *rand.Rand
)

func getRequest(request *Request) *Request {

	if clientWeigh == nil || clientWeigh.weighCount <= 0 {

		return request
	}

	n := uint32(r.Int31n(int32(clientWeigh.weighCount)))

	var (
		count uint32
	)

	for _, value := range clientWeigh.list {
		if count >= n {
			// value.req.Print()
			return value.req
		}
		count = count + value.weights
	}

	panic("getRequest err")

	return nil
}

func getRequestList(request *Request, path string) (clients []*Request) {

	clients = GetRequestListFromFile(path)

	if clients == nil || len(clients) <= 0 {

		return []*Request{request}
	}

	return clients
}

// GetRequestListFromFile is get request from curl file 文件路径为空， 则返回 nil
func GetRequestListFromFile(path string) (clients []*Request) {
	clients = make([]*Request, 0)
	if path == "" {
		return
	}
	curls, err := utils.ParseTheFileC(path)
	if err != nil {
		return
	}
	for _, v := range curls {

		clients = append(clients, &Request{
			Url:     v.GetUrl(),
			Method:  v.GetMethod(),
			Headers: v.GetHeaders(),
			Body:    v.GetBody(),
			Timeout: 30 * time.Second,
			Verify:  VerifyStr,
			Debug:   Debug,
			Form:    FormTypeHttp,
		})
	}
	return
}
