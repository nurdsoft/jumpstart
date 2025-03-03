package pkg

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func runCommands(workdir string, commands Commands) error {
	var err error

	if len(commands.Windows) > 0 {
		err = runCommandsOnWindows(workdir, commands)
		if err != nil {
			return err
		}
	}
	if len(commands.UnixLike) > 0 {
		err = runCommandsOnUnixLike(workdir, commands)
		if err != nil {
			return err
		}
	} else {
		// Run all commands if no OS specific commands are specified
		for _, command := range commands.All {
			cag := strings.Split(command, " ")
			cmd := exec.Command(cag[0], cag[1:]...)
			cmd.Dir = workdir
			if err = cmd.Run(); err != nil {
				return fmt.Errorf("error running command '%s': %w", command, err)
			}
		}
	}
	return nil
}

func runCommandsOnWindows(workdir string, commands Commands) error {
	if runtime.GOOS == "windows" {
		for _, command := range commands.Windows {
			cag := strings.Split(command, " ")
			cmd := exec.Command(cag[0], cag[1:]...)
			cmd.Dir = workdir
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("error running windows command '%s': %w", command, err)
			}
		}
	}
	return nil
}

func runCommandsOnUnixLike(workdir string, commands Commands) error {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		for _, command := range commands.UnixLike {
			cag := strings.Split(command, " ")
			cmd := exec.Command(cag[0], cag[1:]...)
			cmd.Dir = workdir
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("error running unix command '%s': %w", command, err)
			}
		}
	}
	return nil
}
