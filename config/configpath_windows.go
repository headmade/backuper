// +build windows

package config

import "path/filepath"

func configPath() string {
  return filepath.Join("C:", "etc", "gobackuper", "agent.json")
}
