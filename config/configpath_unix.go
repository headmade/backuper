// +build linux darwin freebsd openbsd netbsd dragonfly

package config

import "path/filepath"

func configPath() string {
  return filepath.Join("/", "etc", "gobackuper", "agent.json")
}
