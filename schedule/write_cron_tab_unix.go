// +build linux darwin freebsd openbsd netbsd dragonfly

package schedule

import (
  "fmt"

  "github.com/headmade/backuper/hmutil"
)

func (m *Manager) writeCrontab(schedule string, cmd string) error {
  taskFormat := `crontab -l\
    | ( grep -v 'gobackuper %s' ; echo '%s %s %s gobackuper %s >> /var/log/gobackuper_cron.log 2>&1' )\
    | crontab`
  task := fmt.Sprintf(taskFormat, cmd, schedule, CRON_PATH, CRON_GOTRACEBACK, cmd)
  _, err := hmutil.System(task)
  return err
}
