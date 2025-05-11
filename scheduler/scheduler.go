package scheduler

import (
	// "github.com/adhocore/gronx"
	"github.com/robfig/cron"
)

var CronSheduler = cron.New()

func Schedule(cronString string, job func()) error {
	err := CronSheduler.AddFunc(cronString, job)
	if err != nil {
		return err
	}
	CronSheduler.Start()

	return nil
}
