package gocron

import (
	"strings"
	"fmt"
	"strconv"
	"time"
	"log"
	"os"
)
var CronLog *log.Logger

type CronTicker struct {
	cronExp		string

	secondExp	string
	minuteExp	string
	hourExp		string

	lastTick	int64
}

func init() {
	CronLog = log.New(os.Stdout, "[gocron] ", log.LstdFlags)
}

// 格式: [second minute hour]
// * * * 每秒执行
// 1 * * 每分钟第一秒执行
// */10 * * 每10秒执行
// 5/10 * * 每分钟从第5秒开始, 每10秒执行
// 1 1 23 每天23:01:01执行
func NewCronTicker(exp string) (*CronTicker, error) {
	terms := strings.Split(exp, " ")
	if len(terms) != 3 {
		return nil, fmt.Errorf("invalid cron expression: %s", exp)
	}

	// 合法性验证
	for _, term := range terms {
		if strings.Contains(term, "/") {
			token := strings.Split(term, "/")
			if len(token) != 2 {
				return nil, fmt.Errorf("invalid cron expression: %s", term)
			}

			if valid := isTermValid(token[0]) && isTermValid(token[1]); !valid {
				return nil, fmt.Errorf("invalid cron expression: %s", term)
			}

			break
		}

		if valid := isTermValid(term); !valid {
			return nil, fmt.Errorf("invalid cron expression: %s", term)
		}
	}

	return &CronTicker{
		cronExp: exp,
		secondExp: terms[0],
		minuteExp: terms[1],
		hourExp: terms[2],
	}, nil
}

// 等待下一次执行.
// 该方法会堵塞, 直到下一次执行时间点到来时返回
func (tick *CronTicker) Tick() {
	for {
		tick.waitUntilNextSecond()

		now := time.Now()

		// 匹配秒
		allMatch := tick.match(tick.secondExp, now, 0) && tick.match(tick.minuteExp, now, 1) && tick.match(tick.hourExp, now, 2)

		if allMatch {
			tick.lastTick = now.Unix()
			return
		}

	}
}

func (tick *CronTicker) match(exp string, now time.Time, position int) bool {
	// 通配符
	if exp == "*" {
		return true
	}

	// 是单个数字
	if !strings.Contains(exp, "/") {
		targetTime, err := strconv.Atoi(exp)
		if nil != err {
			CronLog.Println(err)
			return false
		}

		if 0 == position {
			return now.Second() == targetTime
		} else if 1 == position {
			return now.Minute() == targetTime
		} else {
			return now.Hour() == targetTime
		}
	}

	// 是"/"分隔格式
	terms := strings.Split(exp, "/")
	if terms[0] == "*" {
		interval, err := strconv.ParseInt(terms[1], 10, 64)
		if nil != err {
			CronLog.Println(err)
			return false
		}

		diff := now.Unix() - tick.lastTick
		if 0 == position {
			// second
			if diff >= interval * 1 {
				return true
			}
		} else if 1 == position {
			// min
			if diff >= interval * 60 {
				return true
			}
		} else {
			if diff >= interval * 60 * 60 {
				return true
			}
		}

		return false
	} else {
		interval, err := strconv.Atoi(terms[1])
		if nil != err {
			CronLog.Println(err)
			return false
		}


		startFrom, err := strconv.Atoi(terms[0])
		if nil != err {
			CronLog.Println(err)
			return false
		}

		if 0 == position {
			// second
			return now.Second() >= startFrom && (now.Second() - startFrom) % interval == 0
		} else if 1 == position {
			// min
			return now.Minute() >= startFrom && (now.Minute() - startFrom) % interval == 0
		} else {
			return now.Hour() >= startFrom && (now.Hour() - startFrom) % interval == 0
		}

		return false
	}

	CronLog.Println("invalid exp")
	return false
}

func (tick *CronTicker) waitUntilNextSecond() {
	// 计算下一秒的时间
	now := time.Now()
	_1s, _ := time.ParseDuration("+1s")
	nextSecond := now.Add(_1s)

	// 取出下一秒时间的整数秒
	nextRound := time.Date(nextSecond.Year(), nextSecond.Month(), nextSecond.Day(), nextSecond.Hour(), nextSecond.Minute(), nextSecond.Second(), 0, time.Local)

	// 计算时间差
	diff := nextRound.Sub(now)

	// wait
	<- time.After(diff)
}


func isTermValid(str string) bool {
	return isNumber(str) || str == "*"
}

func isNumber(str string) bool {
	_, err := strconv.Atoi(str)
	if nil != err {
		return false
	}

	return true
}