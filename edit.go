package edit // sny.no/tools/edit

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	Editor           = env("EDITOR_B", "$PLAN9/bin/B")
	EditorConnection = env("EDITOR_CONNECTION", "localhost:52670")
	Namespace        = env("NAMESPACE", "")
)

var (
	ErrCannotExec = errors.New("unable to start editor process")
	ErrCircular   = errors.New("EDITOR is circular")
	ErrNotFound   = errors.New("no such editor")
)

// Normalizes the paths of existing files,
// assumes non-existing files are in CWD,
// and skips flags that start with "-".
func Realpaths(args []string) (normargs []string, err error) {
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

func IsProcessCircular(cur, ed string) bool {
	// TODO: check that cur and ed are not the same executable
	return false
}

func env(key, fallbackv string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		v = fallbackv
	}
	return os.ExpandEnv(v)
}
