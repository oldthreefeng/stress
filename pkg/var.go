package pkg

var (
	Concurrency uint64   // 并发请求数
	Number      uint64   // 单个协程的请求总数
	Debug       bool     // 是否开其调试模式
	Path        string   // curl 文件路径
	VerifyStr   string   // 验证方法
	RequestUrl  string   // 压力测试的Url
	Header      []string // 自定义头信息传递给服务器
	Body        string   // HTTP POST 上传的数据
	Compressed bool
)
