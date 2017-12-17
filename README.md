# Gocron

提供秒、分、时纬度的定时功能。



## cron表达式

格式: [second minute hour]

```
* * * 每秒执行
1 * * 每分钟第一秒执行
*/10 * * 每10秒执行
5/10 * * 每分钟从第5秒开始, 每10秒执行
1 1 23 每天23:01:01执行
```

## 使用方法

下载：

```shell
go get github.com/wanghongfei/gocron
```

代码：

```go
tick, _ := NewCronTicker("*/2 * *") // 每2秒tick一次
for {
	tick.Tick() // Tick()方法会堵塞, 直到运行的时间点会返回
	now := time.Now()
	fmt.Printf("%d:%d:%d\n", now.Hour(), now.Minute(), now.Second())
}
```

