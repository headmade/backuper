// +build windows

package schedule

import (
  "fmt"
  "strings"
  "strconv"
  "os"

  "github.com/headmade/backuper/hmutil"
)

func toMinutes(hhmm []string) int {
  hours, _err := strconv.Atoi(hhmm[1])
  minutes, _err := strconv.Atoi(hhmm[0])

  if _err != nil {
    hours = 1
    minutes = 0
  }

  return hours * 60 + minutes
}

func (m *Manager) writeCrontab(schedule string, cmd string) error {
  // taskFormat := "schtasks /Create /TR \"%s %s\" /TN %s /SC MINUTE /MO %v"
  taskFormat := `schtasks /Create /TR "%s %s" /TN %s /SC MINUTE /MO %v`
  time := toMinutes(strings.Split(schedule, " ")[0:2])
  path := "D:\\src\\github.com\\Twizty\\cbackup\\backuper\\backuper.exe"
  taskFormat = fmt.Sprintf(taskFormat, path, cmd, "gobackuper_" + cmd, time)
  f, _ := os.Create("C:\\tmp\\sche.bat")
  f.Write([]byte(taskFormat))
  f.Close()
  res, err := hmutil.System("sche.bat")
  os.Remove("C:\\tmp\\sche.bat")
  return err
}
