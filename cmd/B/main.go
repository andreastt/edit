package main // sny.no/tools/edit/cmd/B

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
	B(flag.Args()...)
}

func B(arg ...string) {
	if os.Getenv("SSH_CONNECTION") != "" {
		if edit.EditorConnection == "" {
			log.Fatal("missing EDITOR_CONNECTiON")
		}
		if err := RB(arg...); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := ExecveB(arg...); err != nil {
			log.Fatal(err)
		}
	}
}

func RB(arg ...string) (err error) {
	// TODO: support unix domain socket
	client, err := rpc.Dial("tcp", edit.EditorConnection)
	if err != nil {
		return err
	}
	normargs, err := edit.Realpaths(arg)
	if err != nil {
		return err
	}
	req := edit.Request{Args: normargs}
	var ex edit.ExitCode
	err = client.Call("Edit.B", req, &ex)
	if err != nil {
		return err
	}
	os.Exit(int(ex))
	panic("never reached")
}

func ExecveB(arg ...string) error {
	b, err := exec.LookPath(edit.Editor)
	if err != nil {
		return err
	}
	cur, err := os.Executable()
	if err != nil {
		return err
	}
	if edit.IsProcessCircular(cur, b) {
		return edit.ErrCircular
	}
	args := append([]string{b}, arg...)
	if err := syscall.Exec(b, args, os.Environ()); err != nil {
		return err
	}
	panic("never reached")
}
