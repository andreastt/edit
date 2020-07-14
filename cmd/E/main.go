package main // sny.no/tools/edit/cmd/E

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"

	"sny.no/tools/edit"
)

const (
	EX_USAGE  = 64
	EX_NOHOST = 68
)

var (
	sshconn = env("SSH_CONNECTION", "")
	edconn  = env("EDITOR_CONNECTION", "localhost:52670")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [+line] file...\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(EX_USAGE)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	flag.Usage = usage
	flag.Parse()
	E(flag.Args()...)
}

func E(arg ...string) {
	if sshconn != "" {
		if len(edconn) < 0 {
			log.Fatal("missing EDITOR_CONNECTION")
		}
		ex, err := RE(arg...)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(int(ex))
	} else {
		if err := edit.ExecveE(arg...); err != nil {
			log.Fatal(err)
		}
	}
}

func RE(arg ...string) (ex edit.ExitCode, err error) {
	client, err := rpc.Dial("tcp", edconn)
	if err != nil {
		return EX_NOHOST, err
	}
	req := edit.Request{
		Args: realpaths(arg),
		//Environ: []string{fmt.Sprintf("NAMESPACE=%s", ns)},
	}
	err = client.Call("Edit.E", req, &ex)
	return ex, err
}

// Normalizes the paths of existing files,
// assumes non-existing files are in CWD,
// and skips flags that start with "-".
func realpaths(args []string) (normargs []string) {
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
			log.Fatal(err)
		}
	}
	return normargs
}

func env(key, fallbackv string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallbackv
	}
	return v
}
