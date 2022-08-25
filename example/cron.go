package example

import (
	"fmt"
	cron "github.com/robfig/cron/v3"
	"time"
)

func TestCron() {
	ticker := time.NewTimer(time.Minute * 1) // 每分钟执行一次
	for range ticker.C {
		fmt.Println("task...")
	}
}

func TestCron1() {
	crontab := cron.New()
	//*    *    *    *   *   *
	//秒   分   时   日  月   年
	//	　每隔5秒执行一次：*/5 * * * * ?
	//	每隔1分钟执行一次：0 */1 * * * ?
	//	每天23点执行一次：0 0 23 * * ?
	//	每天凌晨1点执行一次：0 0 1 * * ?
	//	每月1号凌晨1点执行一次：0 0 1 1 * ?
	//	在26分、29分、33分执行一次：0 26,29,33 * * * ?
	//	每天的0点、13点、18点、21点都执行一次：0 0 0,13,18,21 * * ?
	crontab.AddFunc("* * * * * *", func() {
		println("task")
	})
	crontab.AddFunc("*/3 * * * * *", func() { // 每3s一次
		// do something
	})
	crontab.Start()
	defer crontab.Stop()
	select {} //阻塞
}
