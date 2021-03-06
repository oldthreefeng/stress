package pkg

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// 处理gzip压缩
func GetZipData(response *http.Response) (body []byte, err error) {
	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		defer func() {
			reader.Close()
		}()
	default:
		reader = response.Body
	}

	body, err = ioutil.ReadAll(reader)

	return
}

// HttpStatusCode is 通过Http状态码判断是否请求成功
func HttpStatusCode(request *Request, response *http.Response) (code int, isSucceed bool) {

	defer response.Body.Close()
	code = response.StatusCode
	if code == http.StatusOK {
		isSucceed = true
	}

	// 开启调试模式
	if request.GetDebug() {
		body, err := GetZipData(response)
		fmt.Printf("请求结果 httpCode:%d body:%s err:%v \n", response.StatusCode, string(body), err)

	}

	return
}

/***************************  返回值为json  ********************************/

// 返回数据结构体
type ResponseJson struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// HttpJson is 通过返回的Body 判断
// 返回示例: {"code":200,"msg":"Success","data":{}}
// code 默认将http code作为返回码，http code 为200时 取body中的返回code
func HttpJson(request *Request, response *http.Response) (code int, isSucceed bool) {

	defer response.Body.Close()
	code = response.StatusCode
	if code == http.StatusOK {

		body, err := GetZipData(response)
		if err != nil {
			code = ParseError
			fmt.Printf("请求结果 ioutil.ReadAll err:%v", err)
		} else {
			responseJson := &ResponseJson{}
			err = json.Unmarshal(body, responseJson)
			if err != nil {
				code = ParseError
				fmt.Printf("请求结果 json.Unmarshal err:%v", err)
			} else {

				code = responseJson.Code

				// body 中code返回200为返回数据成功
				if responseJson.Code == 200 {
					isSucceed = true
				}
			}
		}

		// 开启调试模式
		if request.GetDebug() {
			fmt.Printf("请求结果 httpCode:%d body:%s err:%v \n", response.StatusCode, string(body), err)
		}
	}

	return
}

// MallVersion is for mall.youpenglai.com/apis/version
type MallVersion struct {
	Status     int    `json:"status"`
	Result     string `json:"result"`
	TraceId    string `json:"traceId"`
	DevMessage string `json:"devMessage"`
	Message    string `json:"message"`
	Url        string `json:"url"`
}

// MallVersionJson 通过返回的Body 判断
func MallVersionJson(request *Request, response *http.Response) (code int, isSucceed bool) {
	defer response.Body.Close()
	code = response.StatusCode
	if code == http.StatusOK {
		body, err := GetZipData(response)
		if err != nil {
			code = ParseError
			fmt.Printf("请求结果 ioutil.ReadAll err:%v", err)
		} else {
			responseJson := &MallVersion{}
			err = json.Unmarshal(body, responseJson)
			if err != nil {
				code = ParseError
				fmt.Printf("请求结果 json.Unmarshal err:%v", err)
			} else {
				code = responseJson.Status
				// body 中code返回200为返回数据成功
				if responseJson.Status == 200 {
					isSucceed = true
				}
			}
		}
		if request.GetDebug() {
			fmt.Printf("请求结果 httpCode:%d body:%s err:%v \n", response.StatusCode, string(body), err)
		}
	}
	return
}
