package pkg

import (
	"fmt"
	"github.com/oldthreefeng/stress/utils"
	"strings"
	"sync"
	"time"
)

var (
	// 输出统计数据的时间
	exportStatisticsTime = 1 * time.Second
)

// ReceivingResults is 接收结果并处理
// 统计的时间都是纳秒，显示的时间 都是毫秒
// concurrent 并发数
func ReceivingResults(concurrent int, ch <-chan *RequestResults, wg *sync.WaitGroup) {

	defer func() {
		wg.Done()
	}()

	var (
		stopChan = make(chan bool)
	)

	// 时间
	var (
		processingTime uint64 // 处理总时间
		requestTime    uint64 // 请求总时间
		maxTime        uint64 // 最大时长
		minTime        uint64 // 最小时长
		successNum     uint64 // 成功处理数，code为0
		failureNum     uint64 // 处理失败数，code不为0
		chanIdLen      int    // 并发数
		chanIds        = utils.New(Concurrency)
	)

	statTime := uint64(time.Now().UnixNano())

	// 错误码/错误个数
	var errCode = utils.New()

	// 定时输出一次计算结果
	ticker := time.NewTicker(exportStatisticsTime)
	go func() {
		for {
			select {
			case <-ticker.C:
				endTime := uint64(time.Now().UnixNano())
				requestTime = endTime - statTime
				go calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, chanIdLen, errCode)
			case <-stopChan:
				// 处理完成

				return
			}
		}
	}()

	header()

	for data := range ch {
		// fmt.Println("处理一条数据", data.Id, data.Time, data.IsSucceed, data.ErrCode)
		processingTime = processingTime + data.Time

		if maxTime <= data.Time {
			maxTime = data.Time
		}

		if minTime == 0 {
			minTime = data.Time
		} else if minTime > data.Time {
			minTime = data.Time
		}

		// 是否请求成功
		if data.IsSucceed == true {
			successNum = successNum + 1
		} else {
			failureNum = failureNum + 1
		}

		// 统计错误码
		strErrCode := fmt.Sprintf("%d", data.ErrCode)
		if value, ok := errCode.Get(strErrCode); ok {
			errCode.Set(strErrCode, value.(int)+1)
		} else {
			errCode.Set(strErrCode, 1)
		}
		strChanId := fmt.Sprintf("%d", data.ChanId)
		if _, ok := chanIds.Get(strChanId); !ok {
			chanIds.Set(strChanId, true)
			chanIdLen = chanIds.Count()
		}

	}

	// 数据全部接受完成，停止定时输出统计数据
	stopChan <- true

	endTime := uint64(time.Now().UnixNano())
	requestTime = endTime - statTime

	calculateData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, chanIdLen, errCode)

	fmt.Printf("\n\n")

	fmt.Println("*************************  结果 stat  ****************************")
	fmt.Println("处理协程数量:", concurrent)
	// fmt.Println("处理协程数量:", concurrent, "程序处理总时长:", fmt.Sprintf("%.3f", float64(processingTime/concurrent)/1e9), "秒")
	fmt.Println("请求总数（并发数*请求数 -c * -n）:", successNum+failureNum, "总请求时间:", fmt.Sprintf("%.3f", float64(requestTime)/1e9),
		"秒", "successNum:", successNum, "failureNum:", failureNum)

	fmt.Println("*************************  结果 end   ****************************")

	fmt.Printf("\n\n")
}

// 计算数据
func calculateData(concurrent int, processingTime, requestTime, maxTime, minTime, successNum, failureNum uint64, chanIdLen int, errCode utils.ConcurrentMap) {
	if processingTime == 0 {
		processingTime = 1
	}

	var (
		qps              float64
		averageTime      float64
		maxTimeFloat     float64
		minTimeFloat     float64
		requestTimeFloat float64
	)

	// 平均 每个协程成功数*总协程数据/总耗时 (每秒)
	if processingTime != 0 {
		qps = float64(successNum*1e9*uint64(concurrent)) / float64(processingTime)
	}

	// 平均时长 总耗时/总请求数/并发数 纳秒=>毫秒
	if successNum != 0 && concurrent != 0 {
		averageTime = float64(processingTime) / float64(successNum*1e6*uint64(concurrent))
	}

	// 纳秒=>毫秒
	maxTimeFloat = float64(maxTime) / 1e6
	minTimeFloat = float64(minTime) / 1e6
	requestTimeFloat = float64(requestTime) / 1e9

	// 打印的时长都为毫秒
	// result := fmt.Sprintf("请求总数:%8d|successNum:%8d|failureNum:%8d|qps:%9.3f|maxTime:%9.3f|minTime:%9.3f|平均时长:%9.3f|errCode:%v", successNum+failureNum, successNum, failureNum, qps, maxTimeFloat, minTimeFloat, averageTime, errCode)
	// fmt.Println(result)
	table(successNum, failureNum, errCode, qps, averageTime, maxTimeFloat, minTimeFloat, requestTimeFloat, chanIdLen)
}

// 打印表头信息
func header() {
	fmt.Printf("\n\n")
	// 打印的时长都为毫秒 总请数
	fmt.Println("─────┬───────┬───────┬───────┬────────┬────────┬────────┬────────┬────────")
	result := fmt.Sprintf(" 耗时│ 并发数│ 成功数│ 失败数│   qps  │最长耗时│最短耗时│平均耗时│ 错误码")
	fmt.Println(result)
	// result = fmt.Sprintf("耗时(s)  │总请求数│成功数│失败数│QPS│最长耗时│最短耗时│平均耗时│错误码")
	// fmt.Println(result)
	fmt.Println("─────┼───────┼───────┼───────┼────────┼────────┼────────┼────────┼────────")

	return
}

// 打印表格
func table(successNum, failureNum uint64, errCode utils.ConcurrentMap, qps, averageTime, maxTimeFloat, minTimeFloat, requestTimeFloat float64, chanIdLen int) {
	// 打印的时长都为毫秒
	result := fmt.Sprintf("%4.0fs│%7d│%7d│%7d│%8.2f│%8.2f│%8.2f│%8.2f│%v", requestTimeFloat, chanIdLen, successNum, failureNum, qps, maxTimeFloat, minTimeFloat, averageTime, printMap(errCode))
	fmt.Println(result)

	return
}

// 输出错误码、次数 节约字符(终端一行字符大小有限)
func printMap(errCode utils.ConcurrentMap) (mapStr string) {

	var (
		mapArr []string
	)
	for key, value := range errCode.Items() {
		mapArr = append(mapArr, fmt.Sprintf("%s:%d", key, value))
	}

	mapStr = strings.Join(mapArr, ";")

	return
}
