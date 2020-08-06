# stress

go 实现的压测工具，每个用户用一个协程的方式模拟，最大限度的利用 CPU 资源

### 1.1 安装

使用make安装. 建议使用make安装.

```shell
$ git clone https://github.com/oldthreefeng/stress
$ make linux
```

使用`golang build`

```shell
$ git clone https://github.com/oldthreefeng/stress
$ go build -o stress main.go
```

使用脚本项目脚本构建 `build.sh`

```shell 
$ git clone https://github.com/oldthreefeng/stress
$ sh build.sh v1.0.0
```

### 1.2 项目体验


参数说明:

`-c` 表示并发数 默认为 1

`-n` 每个并发执行请求的次数， 默认为 1 ，总请求的次数 = 并发数 `*` 每个并发执行请求的次数

`-u` 需要压测的地址 example: http://www.baidu.com/

```shell

# 运行 以linux为示例
./stress -c 1 -n 100 -u https://www.baidu.com/

```

- 压测结果展示

执行以后，终端每秒钟都会输出一次结果，压测完成以后输出执行的压测结果

压测结果展示:

```

─────┬───────┬───────┬───────┬────────┬────────┬────────┬────────┬────────
 耗时│ 并发数 │ 成功数│ 失败数 │   qps  │最长耗时 │最短耗时│平均耗时 │ 错误码
─────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────
   1s│      1│      8│      0│    8.09│  133.16│  110.98│  123.56│200:8
   2s│      1│     15│      0│    8.02│  138.74│  110.98│  124.61│200:15
   3s│      1│     23│      0│    7.80│  220.43│  110.98│  128.18│200:23
   4s│      1│     31│      0│    7.83│  220.43│  110.23│  127.67│200:31
   5s│      1│     39│      0│    7.81│  220.43│  110.23│  128.03│200:39
   6s│      1│     46│      0│    7.72│  220.43│  110.23│  129.59│200:46
   7s│      1│     54│      0│    7.79│  220.43│  110.23│  128.42│200:54
   8s│      1│     62│      0│    7.81│  220.43│  110.23│  128.09│200:62
   9s│      1│     70│      0│    7.79│  220.43│  110.23│  128.33│200:70
  10s│      1│     78│      0│    7.82│  220.43│  106.47│  127.85│200:78
  11s│      1│     84│      0│    7.64│  371.02│  106.47│  130.96│200:84
  12s│      1│     91│      0│    7.63│  371.02│  106.47│  131.02│200:91
  13s│      1│     99│      0│    7.66│  371.02│  106.47│  130.54│200:99
  13s│      1│    100│      0│    7.66│  371.02│  106.47│  130.52│200:100


*************************  结果 stat  ****************************
处理协程数量: 1
请求总数: 100 总请求时间: 13.055 秒 successNum: 100 failureNum: 0
*************************  结果 end   ****************************

```

参数解释:

**耗时**: 程序运行耗时。程序每秒钟输出一次压测结果

**并发数**: 并发数，启动的协程数

**成功数**: 压测中，请求成功的数量

**失败数**: 压测中，请求失败的数量

**qps**: 当前压测的QPS(每秒钟处理请求数量)

**最长耗时**: 压测中，单个请求最长的响应时长

**最短耗时**: 压测中，单个请求最短的响应时长

**平均耗时**: 压测中，单个请求平均的响应时长

**错误码**: 压测中，接口返回的 code码:返回次数的集合


### 2.1 压测是什么

压测，即压力测试，是确立系统稳定性的一种测试方法，通常在系统正常运作范围之外进行，以考察其功能极限和隐患。

主要检测服务器的承受能力，包括用户承受能力（多少用户同时玩基本不影响质量）、流量承受等。

- 压测的目的就是通过压测(模拟真实用户的行为)，测算出机器的性能(单台机器的QPS)，从而推算出系统在承受指定用户数(100W)时，需要多少机器能支撑得住
- 压测是在上线前为了应对未来可能达到的用户数量的一次预估(提前演练)，压测以后通过优化程序的性能或准备充足的机器，来保证用户的体验。

| 压测类型 |   解释  |
| :----   | :---- |
| 压力测试(Stress Testing)          |  也称之为强度测试，测试一个系统的最大抗压能力，在强负载(大数据、高并发)的情况下，测试系统所能承受的最大压力，预估系统的瓶颈    |
| 并发测试(Concurrency Testing)     |  通过模拟很多用户同一时刻访问系统或对系统某一个功能进行操作，来测试系统的性能，从中发现问题(并发读写、线程控制、资源争抢)      |
| 耐久性测试(Configuration Testing) |  通过对系统在大负荷的条件下长时间运行，测试系统、机器的长时间运行下的状况,从中发现问题(内存泄漏、数据库连接池不释放、资源不回收)     |

### 2.2 如何计算压测指标

- 压测我们需要有目的性的压测，这次压测我们需要达到什么目标(如:单台机器的性能为 100QPS?网站能同时满足100W人同时在线)
- 可以通过以下计算方法来进行计算:
- 压测原则:每天80%的访问量集中在20%的时间里，这20%的时间就叫做峰值
- 公式: ( 总PV数`*`80% ) / ( 每天的秒数`*`20% ) = 峰值时间每秒钟请求数(QPS)
- 机器: 峰值时间每秒钟请求数(QPS) / 单台机器的QPS = 需要的机器的数量

- 假设:网站每天的用户数(100W)，每天的用户的访问量约为3000W PV，这台机器的需要多少QPS?
> ( 30000000\*0.8 ) / (86400 * 0.2) ≈ 1389 (QPS)
>


### 3.1 用法

```shell
$ ./stress -h
stress is a test cli for http and websocket stress written by golang, 
go 实现的压测工具，每个用户用一个协程的方式模拟，最大限度的利用 CPU 资源

Usage:
  stress [flags]
  stress [command]

Examples:
        
        # stress curl file to test 
        stress -f utils/curl.txt

        # stress curl file read from stdin 
        cat utils/curl.txt | stress -f -

        # stress concurrency 10 & 10 times
        stress -c 10 -n 10 -f  utils/curl.txt

        # stress cli url
        stress -c 10 -n 100 -u https://www.baidu.com


Available Commands:
  help        Help about any command
  version     stress version

Flags:
  -c, --concurrency uint    并发数 (default 1)
      --config string       config file for stress (default is $HOME/.stress.yaml)
      --data string         http post data
  -d, --debug               debug 模式
  -H, --header strings      http post data
  -h, --help                help for stress
  -n, --number uint         单协程的请求数 (default 1)
  -f, --path string         read curl file to build test
  -u, --requestUrl string   curl文件路径
  -t, --toggle              Help message for toggle
  -v, --verify string        verify 验证方法 在server/verify中 http 支持:statusCode、json webSocket支持:json (default "statusCode")

Use "stress [command] --help" for more information about a command.

```


- 使用 curl文件进行压测

curl是Linux在命令行下的工作的文件传输工具，是一款很强大的http命令行工具。

使用curl文件可以压测使用非GET的请求，支持设置http请求的 method、cookies、header、body等参数


**I:** chrome 浏览器生成 curl文件，打开开发者模式(快捷键F12)，如图所示，生成 curl 在终端执行命令
![chrome cURL](https://img.mukewang.com/5d60eddd0001f4b016661114.png)

**II:** postman 生成 curl 命令
![postman cURL](https://img.mukewang.com/5ed79b590001837120581530.png)

生成内容粘贴到项目目录下的**utils/curl.txt**文件中，执行下面命令就可以从curl.txt文件中读取需要压测的内容进行压测了

- 支持多步压力测试

目前使用的方法是按行来分割每个请求. 


```shell
$ cat utils/curl.txt
curl 'https://mall.youpenglai.com/apis/version' -H 'User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:79.0) Gecko/20100101 Firefox/79.0' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8' -H 'Accept-Language: zh-CN,en-US;q=0.7,en;q=0.3' --compressed -H 'Connection: keep-alive' -H 'Upgrade-Insecure-Requests: 1' -H 'Cache-Control: max-age=0' -H 'TE: Trailers'
curl 'https://www.youpenglai.com/v2/sys/server' -H 'User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:79.0) Gecko/20100101 Firefox/79.0' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8' -H 'Accept-Language: zh-CN,en-US;q=0.7,en;q=0.3' --compressed -H 'Connection: keep-alive' -H 'Cookie: _ga=GA1.2.462595907.1596072242; Hm_lvt_16789552d8e4ee2108c29cc8197bfe19=1596077258; _gid=GA1.2.1014297080.1596441772; JSESSIONID=88EE56D0223D34DD83EB301962B121D1' -H 'Upgrade-Insecure-Requests: 1' -H 'Cache-Control: max-age=0'
```

支持从`stdin`读取`curl`请求. 从`stdin`读取的内容会保存在`/tmp/curl.tmp`

```shell script
$ cat utils/curl.txt  | ./stress -f -
```

命令行示例. 

```shell
$ ./stress -c 1 -n 1 -d -u 'https://page.aliyun.com/delivery/plan/list' \
  -H 'authority: page.aliyun.com' \
  -H 'accept: application/json, text/plain, */*' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -H 'origin: https://cn.aliyun.com' \
  -H 'sec-fetch-site: same-site' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-dest: empty' \
  -H 'referer: https://cn.aliyun.com/' \
  -H 'accept-language: zh-CN,zh;q=0.9' \
  -H 'cookie: aliyun_choice=CN; JSESSIONID=J8866281-CKCFJ4BUZ7GDO9V89YBW1-KJ3J5V9K-GYUW7; maliyun_temporary_console0=1AbLByOMHeZe3G41KYd5WWZvrM%2BGErkaLcWfBbgveKA9ifboArprPASvFUUfhwHtt44qsDwVqMk8Wkdr1F5LccYk2mPCZJiXb0q%2Bllj5u3SQGQurtyPqnG489y%2FkoA%2FEvOwsXJTvXTFQPK%2BGJD4FJg%3D%3D; cna=L3Q5F8cHDGgCAXL3r8fEZtdU; isg=BFNThsmSCcgX-sUcc5Jo2s2T4tF9COfKYi8g9wVwr3KphHMmjdh3GrHFvPTqJD_C; l=eBaceXLnQGBjstRJBOfwPurza77OSIRAguPzaNbMiT5POw1B5WAlWZbqyNY6C3GVh6lwR37EODnaBeYBc3K-nxvOu9eFfGMmn' \
  --data 'adPlanQueryParam=%7B%22adZone%22%3A%7B%22positionList%22%3A%5B%7B%22positionId%22%3A83%7D%5D%7D%2C%22requestId%22%3A%2217958651-f205-44c7-ad5d-f8af92a6217a%22%7D'
  --compressed
```

like curl cmd , use `--data-raw` is the same with `--data`,  when use `--compressed`,  only support command line and `gzip`. 
for curl file. the Header is in it. use `--compressed`, just  do this `req.Header.Add("Accept-Encoding", "gzip")`. 

```shell
$ ./stress -c 1 -n 1 -d -u 'https://page.aliyun.com/delivery/plan/list' \
  -H 'authority: page.aliyun.com' \
  -H 'accept: application/json, text/plain, */*' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -H 'origin: https://cn.aliyun.com' \
  -H 'sec-fetch-site: same-site' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-dest: empty' \
  -H 'referer: https://cn.aliyun.com/' \
  -H 'accept-language: zh-CN,zh;q=0.9' \
  -H 'cookie: aliyun_choice=CN; JSESSIONID=J8866281-CKCFJ4BUZ7GDO9V89YBW1-KJ3J5V9K-GYUW7; maliyun_temporary_console0=1AbLByOMHeZe3G41KYd5WWZvrM%2BGErkaLcWfBbgveKA9ifboArprPASvFUUfhwHtt44qsDwVqMk8Wkdr1F5LccYk2mPCZJiXb0q%2Bllj5u3SQGQurtyPqnG489y%2FkoA%2FEvOwsXJTvXTFQPK%2BGJD4FJg%3D%3D; cna=L3Q5F8cHDGgCAXL3r8fEZtdU; isg=BFNThsmSCcgX-sUcc5Jo2s2T4tF9COfKYi8g9wVwr3KphHMmjdh3GrHFvPTqJD_C; l=eBaceXLnQGBjstRJBOfwPurza77OSIRAguPzaNbMiT5POw1B5WAlWZbqyNY6C3GVh6lwR37EODnaBeYBc3K-nxvOu9eFfGMmn' \
  --data-raw 'adPlanQueryParam=%7B%22adZone%22%3A%7B%22positionList%22%3A%5B%7B%22positionId%22%3A83%7D%5D%7D%2C%22requestId%22%3A%2217958651-f205-44c7-ad5d-f8af92a6217a%22%7D'
  --compressed
```

- 项目结构

```
├── build.sh
├── cmd
│   ├── dispose.go
│   ├── root.go
│   └── version.go
├── main.go
├── Makefile
├── pkg
│   ├── http_client.go
│   ├── websocket_client.go
│   ├── http.go
│   ├── request.go
│   ├── statistics.go
│   ├── var.go
│   ├── verify_http.go
│   ├── verify_websocket.go
│   └── websocket.go
├── README.md
└── utils
    ├── curl.txt
    ├── utils.go
    └── utils_test.go
```

### 4.参考文献
    
[性能测试工具](https://testerhome.com/topics/17068)

[性能测试常见名词解释](https://blog.csdn.net/r455678/article/details/53063989)

[性能测试名词解释](https://codeigniter.org.cn/forums/blog-39678-2456.html)

[PV、TPS、QPS是怎么计算出来的？](https://www.zhihu.com/question/21556347)

[超实用压力测试工具－ab工具](https://www.jianshu.com/p/43d04d8baaf7)