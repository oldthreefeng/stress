package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// Curl store curl file use for go map struct
type Curl struct {
	Data map[string][]string
}

func (c *Curl) getDataValue(k []string) []string {
	var v = make([]string, 0)
	for _, key := range k {
		var ok bool
		v, ok = c.Data[key]
		if ok {
			break
		}
	}
	return v
}

func ParseTheFile(path string) (curl *Curl, err error) {
	if path == "" {
		err = errors.New("路径不能为空")

		return
	}
	file, err := os.Open(path)
	if err != nil {
		err = errors.New("打开文件失败:" + err.Error())

		return
	}

	defer func() {
		file.Close()
	}()

	dataBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	return NewCurl(string(dataBytes)), nil
}

// ParseTheFileC is read  multi curl file to Curl struct
func ParseTheFileC(path string) (curls []*Curl, err error) {

	if path == "" {
		err = errors.New("路径不能为空")

		return
	}
	file, err := os.Open(path)
	if err != nil {
		err = errors.New("打开文件失败:" + err.Error())

		return
	}

	defer func() {
		file.Close()
	}()

	//dataBytes, err := ioutil.ReadAll(file)
	br := bufio.NewReader(file)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		curl := NewCurl(string(a))
		curls = append(curls, curl)
	}

	// for key, value := range curl.Data {
	// 	fmt.Println("key:", key, "value:", value)
	// }

	return
}

// NewCurl is read string to Curl
func NewCurl(data string) (curl *Curl) {
	curl = &Curl{
		Data: make(map[string][]string),
	}

	for len(data) > 0 {
		if strings.HasPrefix(data, "curl") {
			data = data[5:]
		}

		data = strings.TrimSpace(data)
		var (
			key   string
			value string
		)
		index := strings.Index(data, " ")
		if index <= 0 {
			break
		}
		key = strings.TrimSpace(data[:index])
		data = data[index+1:]
		data = strings.TrimSpace(data)

		// url
		if !strings.HasPrefix(key, "-") {
			key = strings.Trim(key, "'")
			curl.Data["curl"] = []string{key}

			// 去除首尾空格
			data = strings.TrimFunc(data, func(r rune) bool {
				if r == ' ' || r == '\\' || r == '\n' {
					return true
				}

				return false
			})
			continue
		}

		if strings.HasPrefix(data, "-") {
			continue
		}

		var (
			endSymbol = " "
		)

		if strings.HasPrefix(data, "'") {
			endSymbol = "'"
			data = data[1:]
		}

		index = strings.Index(data, endSymbol)
		if index <= -1 {
			break
		}
		value = data[:index]
		data = data[index+1:]

		// 去除首尾空格
		data = strings.TrimFunc(data, func(r rune) bool {
			if r == ' ' || r == '\\' || r == '\n' {
				return true
			}

			return false
		})
		curl.Data[key] = append(curl.Data[key], value)

	}
	return
}

// GetMethod is read Method from Curl
func (c *Curl) GetMethod() (method string) {
	method = "GET"

	var (
		postKeys = []string{"--d", "--data", "--data-binary $", "--data-binary"}
	)
	value := c.getDataValue(postKeys)

	if len(value) >= 1 {
		return "POST"
	}

	keys := []string{"-X", "--request"}
	value = c.getDataValue(keys)

	if len(value) <= 0 {

		return
	}

	method = strings.ToUpper(value[0])

	return
}

// GetHeaders is read Headers from Curl
func (c *Curl) GetHeaders() (headers ConcurrentMap) {

	keys := []string{"-H", "--header"}
	value := c.getDataValue(keys)

	for _, v := range value {
		getHeaderValue(v, headers)
	}

	return
}

func getHeaderValue(v string, headers ConcurrentMap) {
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

// GetHeadersStr is get Header returns string
func (c *Curl) GetHeadersStr() string {
	headers := c.GetHeaders()
	header := headers.Items()
	bytes, _ := json.Marshal(&header)

	return string(bytes)
}

// GetBody is get Body from Curl
func (c *Curl) GetBody() (body string) {

	keys := []string{"--data", "-d", "--data-raw", "--data-binary"}
	value := c.getDataValue(keys)

	if len(value) <= 0 {

		return
	}

	// body = strings.NewReader(value[0])
	body = value[0]

	return
}

// GetUrl is get url from Curl
func (c *Curl) GetUrl() (url string) {

	keys := []string{"curl", "--url"}
	value := c.getDataValue(keys)
	if len(value) <= 0 {

		return
	}

	url = value[0]

	return
}

// DiffNano is return start & end time for program
func DiffNano(startTime time.Time) (diff int64) {

	startTimeStamp := startTime.UnixNano()
	endTimeStamp := time.Now().UnixNano()

	diff = endTimeStamp - startTimeStamp

	return
}

// FileExist return true if file exist
func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
