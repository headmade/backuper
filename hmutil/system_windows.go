// +build windows

package hmutil

import (
  "log"
  "os/exec"
)

func System(cmd string) ([]byte, error) {
  log.Println(cmd)
  return exec.Command("cmd", "/c", cmd).CombinedOutput()
}
