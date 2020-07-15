package main // sny.no/tools/edit/cmd/E

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"syscall"

	"sny.no/tools/edit"
)

const EX_USAGE = 64

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
	if os.Getenv("SSH_CONNECTION") != "" {
		if edit.EditorConnection == "" {
			log.Fatal("missing EDITOR_CONNECTION")
		}
		if err := RE(arg...); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := ExecveE(arg...); err != nil {
			log.Fatal(err)
		}
	}
}

func RE(arg ...string) (err error) {
	client, err := rpc.Dial("tcp", edit.EditorConnection)
	if err != nil {
		return err
	}
	normargs, err := edit.Realpaths(arg)
	if err != nil {
		return err
	}
	req := edit.Request{
		Args: normargs,
		//Environ: []string{fmt.Sprintf("NAMESPACE=%s", ns)},
	}
	var ex edit.ExitCode
	err = client.Call("Edit.E", req, &ex)
	if err != nil {
		return err
	}
	os.Exit(int(ex))
	panic("never reached")
}

func ExecveE(arg ...string) error {
	ed, err := exec.LookPath(edit.Editor)
	if err != nil {
		return err
	}

	// TODO: check for circular EDITOR

	args := append([]string{ed}, arg...)
	if err := syscall.Exec(ed, args, os.Environ()); err != nil {
		return err
	}
	panic("never reached")
}
