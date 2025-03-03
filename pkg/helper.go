package pkg

import (
	"io"
	"os"
	"path/filepath"
)

func AbsPath(dest string) string {
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

func copyFile(src, dest string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}
