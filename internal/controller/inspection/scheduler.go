package inspection

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	devopsv1 "github.com/rxg456/auto-inspection-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNextRunTime 使用cron库获取下一次运行时间
func GetNextRunTime(job devopsv1.Job, lastTime *metav1.Time) (time.Time, error) {
	// 创建cron解析器
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	// 解析cron表达式
	schedule, err := parser.Parse(job.Schedule)
	if err != nil {
		return time.Time{}, fmt.Errorf("解析cron表达式失败: %w", err)
	}

	now := time.Now()
	var nextRun time.Time

	if lastTime == nil {
		// 如果没有上次执行时间，立即执行
		nextRun = now
	} else {
		// 计算下一次执行时间
		nextRun = schedule.Next(lastTime.Time)

		// 如果下一次执行时间已经过了，立即执行
		if nextRun.Before(now) {
			nextRun = now
		}
	}

	return nextRun, nil
}

// ShouldRunNow 判断任务是否应该立即执行
func ShouldRunNow(job devopsv1.Job, lastRunTime *metav1.Time) (bool, error) {
	nextRunTime, err := GetNextRunTime(job, lastRunTime)
	if err != nil {
		return false, err
	}

	// 如果下一次执行时间不在未来，则应该立即运行
	return !nextRunTime.After(time.Now()), nil
}
