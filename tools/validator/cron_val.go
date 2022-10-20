package validator

import "github.com/jakecoffman/cron"

// ValCron 校验Cron表达式是否正确
// 如果不正确返回true
// 正确返回false
func ValCron(cronStr string) (isNotCorrect bool) {
	defer func() {
		if err := recover(); err != nil {
			isNotCorrect = true
		}
	}()
	_ = cron.Parse(cronStr)
	return
}
