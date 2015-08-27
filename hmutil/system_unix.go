// +build linux darwin freebsd openbsd netbsd dragonfly

package hmutil

import (
  "log"
  "os/exec"
)

func System(cmd string) ([]byte, error) {
  log.Println(cmd)
  return exec.Command("sh", "-c", cmd).CombinedOutput()
}
