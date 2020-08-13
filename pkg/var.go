package pkg

var (
	Concurrency int   // Concurrency is 并发请求数
	Number      uint64   // Number is 单个协程的请求总数
	Debug       bool     // Debug is 是否开其调试模式
	Path        string   // Path is curl 文件路径
	VerifyStr   string   // VerifyStr is 验证方法
	RequestUrl  string   // RequestUrl is 压力测试的Url
	Header      []string // Header is 自定义头信息传递给服务器
	Body        string   // Body is HTTP POST data 上传的数据
	Compressed  bool     // Compressed is curl zip algorithm
)
