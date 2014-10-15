package schedule

import (
	"fmt"
	"io/ioutil"

	"github.com/headmade/backuper/backuper"
)

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
	task := fmt.Sprintf("%s root /usr/local/bin/gobackuper %s >> /var/log/gobackuper_cron.log 2>&1\n\n", schedule, cmd)
	//return ioutil.WriteFile("/etc/cron.d/backuper_"+cmd, []byte(task), 0600)
	return ioutil.WriteFile("/tmp/gobackuper_"+cmd, []byte(task), 0600)
}
