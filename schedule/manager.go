package schedule

import (
	"fmt"
//	"strconv"

	"github.com/headmade/backuper/backuper"
	"github.com/headmade/backuper/hmutil"
)

const (
	GOTRACEBACK = 1
	CRON_PATH = "PATH=/bin:/usr/bin:/usr/local/bin:$PATH"
)

var CRON_GOTRACEBACK = "" //"CRON_GOTRACEBACK=" + strconv.Itoa(GOTRACEBACK_LEVEL)

type Manager struct {
}

func New() (Manager, error) {
	return Manager{}, nil
}

func (m *Manager) UpdateBackup(period backuper.Period) error {
	return m.writeCrontab(periodToStr(&period), "backup")
}

func (m *Manager) UpdateCheck() error {
	return m.writeCrontab("* * * * *", "check")
}

func periodToStr(period *backuper.Period) string {
	str := "* * * * *"
	return str
}

func (m *Manager) writeCrontab(schedule string, cmd string) error {
	taskFormat := `crontab -l\
		| ( grep -v 'gobackuper %s' ; echo '%s %s %s gobackuper %s >> /var/log/gobackuper_cron.log 2>&1' )\
		| crontab`
	task := fmt.Sprintf(taskFormat, cmd, schedule, CRON_PATH, CRON_GOTRACEBACK, cmd)
	_, err := hmutil.System(task)
	return err
}
