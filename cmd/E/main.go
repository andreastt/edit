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

	"github.com/fsnotify/fsnotify"
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
	start(flag.Args()...)
}

func start(arg ...string) {
	if os.Getenv("SSH_CONNECTION") != "" {
		if edit.EditorConnection == "" {
			log.Fatal("missing EDITOR_CONNECTION")
		}
		if err := RE(arg...); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := E(arg...); err != nil {
			log.Fatal(err)
		}
	}
}

func RE(arg ...string) (err error) {
	// TODO: support unix domain socket
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

// TODO: watch parent directory,
// and if it doesn't exist, watch its parent recursively,
// so we can watch for file creation
func E(arg ...string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					done <- true
				}
			case err := <-watcher.Errors:
				log.Println("watcher:", err)
			}
		}
	}()
	// TODO: support multiple files
	if err := watcher.Add(arg[0]); err != nil {
		return err
	}
	B(arg...)
	<-done
	return nil
}

func B(arg ...string) error {
	b, err := exec.LookPath(edit.Editor)
	if err != nil {
		return err
	}
	cmd := exec.Command(b, arg...)
	if err := cmd.Start(); err != nil {
		return edit.ErrCannotExec
	}
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Fatal("B returned with non-zero exit code:", status)
			}
		}
	}
	return nil
}

/*
// Find nearest parent directory that exists.
func findNearestDirectory(path string) (parent string, err error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	parent = filepath.Dir(abs)
	fmt.Println("testing", parent)
	if _, err := os.Stat(parent); err == nil {
		fmt.Println("exists", parent)
		return filepath.Abs(parent)
	} else if os.IsNotExist(err) {
		fmt.Println("does not exist", parent)
		return findNearestDirectory(parent)
	} else {
		// Schrödinger: file may or may not exist
		panic("schrödinger")
	}
	return
}
*/
