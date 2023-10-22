package pkg

import (
	"os"
	"path/filepath"
	"strings"
)

type SSHKeys struct {
	selected int
	dir      string
	keys     []string
}

func NewSSHKeys() *SSHKeys {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	sshKeysDir := filepath.Join(homedir, ".ssh")
	files, err := os.ReadDir(sshKeysDir)
	if err != nil {
		panic(err)
	}

	keys := make([]string, 0)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if filepath.Ext(f.Name()) == ".pub" {
			if _, err := os.Stat(filepath.Join(sshKeysDir, f.Name())); err != nil {
				continue
			}

			splits := strings.Split(f.Name(), ".")
			key := strings.Join(splits[:len(splits)-1], ".")
			keys = append(keys, key)
		}
	}
	return &SSHKeys{dir: sshKeysDir, keys: keys}
}

func (s *SSHKeys) List() []string {
	return s.keys
}

func (s *SSHKeys) Get() string {
	return filepath.Join(s.dir, s.keys[s.selected])
}

func (s *SSHKeys) SetSelected(i int) {
	s.selected = i
}
