package edit // sny.no/tools/edit

import (
	"os"
	"os/exec"
	"errors"
	"net/rpc"
	"syscall"
	"strings"
	"path/filepath"
)

// TODO: document
const (
	EX_OK         = 0
	EX_USAGE = 64
	EX_NOHOST = 68
	EX_IOERR = 74
	EX_CANNOTEXEC = 126
	EX_NOTFOUND   = 127
	EX_INVALID    = 128
)

var (
	Editor    = env("EDITOR", "ed")
	EditorConnection = env("EDITOR_CONNECTION", "localhost:52670")
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

// TODO: make specific errors so that we don't return specific ExitCodes here
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

func RE(arg ...string) (ex ExitCode, err error) {
	if len(EditorConnection) < 0 {
		return EX_USAGE, errors.New("missing EDITOR_CONNECTION")
	}
	client, err := rpc.Dial("tcp", EditorConnection)
	if err != nil {
		return EX_NOHOST, err
	}
	normargs, err := realpaths(arg)
	if err != nil {
		return EX_IOERR, err
	}
	req := Request{
		Args: normargs,
		//Environ: []string{fmt.Sprintf("NAMESPACE=%s", ns)},
	}
	err = client.Call("Edit.E", req, &ex)
	return ex, err
}

// Normalizes the paths of existing files,
// assumes non-existing files are in CWD,
// and skips flags that start with "-".
func realpaths(args []string) (normargs []string, err error) {
	for _, arg := range args {
		if strings.Index(arg, "-") == 0 {
			normargs = append(normargs, arg)
		} else if _, err := os.Stat(arg); err == nil {
			path, _ := filepath.Abs(arg)
			normargs = append(normargs, path)
		} else if os.IsNotExist(err) {
			cwd, _ := os.Getwd()
			path := filepath.Join(cwd, arg)
			normargs = append(normargs, path)
		} else {
			// Schrödinger: file may or may not exist
			return nil, err
		}
	}
	return normargs, nil
}

func env(key, fallbackv string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallbackv
	}
	return v
}
