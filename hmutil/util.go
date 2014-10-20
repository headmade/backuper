package hmutil

import (
	"os/exec"
	"strings"
)

func System(cmd string) ([]byte, error) {
	return exec.Command("sh", "-c", cmd).CombinedOutput()
}

func ReplaceVars(str string, replacements map[string]string) string {
	for from, to := range replacements {
		str = strings.Replace(str, from, to, -1)
	}
	return str
}

/*
func ErrString(err error) (s *string) {
  if err != nil {
		tmp := err.Error()
		s = &tmp
	}
	return
}
*/
