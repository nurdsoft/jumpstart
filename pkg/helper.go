package pkg

import (
	"os"
	"path/filepath"
)

func absPath(dest string) string {
	if dest[:2] == "~/" {
		homedir, _ := os.UserHomeDir()
		return homedir + dest[1:]
	}

	if dest[0] == '.' {
		abs, _ := filepath.Abs(dest)
		return abs
	}

	return dest
}
