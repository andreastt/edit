package edit // sny.no/tools/edit

import (
	"os"
	"os/exec"
	"syscall"
)

// TODO: document
const (
	EX_OK         = 0
	EX_CANNOTEXEC = 126
	EX_NOTFOUND   = 127
	EX_INVALID    = 128
)

var (
	Editor    = env("EDITOR", "ed")
	Namespace = env("NAMESPACE", "")
)

// Executes EDITOR and replaces the current process with the editor's.
func ExecveE(arg ...string) error {
	ed, err := exec.LookPath(Editor)
	if err != nil {
		return err
	}

	// TODO: check for circular EDITOR

	args := append([]string{ed}, arg...)
	if err := syscall.Exec(ed, args, os.Environ()); err != nil {
		return err
	}
	return nil
}

func E(arg ...string) (ExitCode, error) {
	ed, err := exec.LookPath(Editor)
	if err != nil {
		return EX_NOTFOUND, err
	}
	cmd := exec.Command(ed, arg...)
	if err := cmd.Start(); err != nil {
		return EX_CANNOTEXEC, err
	}
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return ExitCode(status.ExitStatus()), err
			}
		} else {
			return EX_INVALID, err
		}
	}
	return EX_OK, nil
}

func env(key, fallbackv string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallbackv
	}
	return v
}
