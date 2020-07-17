package main // sny.no/tools/edit/cmd/editd

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"sny.no/tools/edit"
)

const EX_USAGE = 64

var (
	network = "tcp"
	addr    = ":52670"
	mountps []mountpoint
)

var (
	afnet   = flag.String("U", "", "use Unix domain socket")
	verbose = flag.Bool("v", false, "log debug messages")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [-v] [<hostname>[:<port>]]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(EX_USAGE)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	flag.Usage = usage
	flag.Parse()

	if len(*afnet) > 0 {
		network = "unix"
		addr = *afnet
	} else if flag.NArg() == 1 {
		addr = flag.Arg(0)
	}

	var err error
	if mountps, err = mountpoints(); err != nil {
		log.Println("unable to list mountpoints:", err)
	}

	so, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer so.Close()
	if so.Addr().Network() == "unix" {
		// cleanup Unix domain socket on exit
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			if err := os.Remove(addr); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		}()
	}

	log.Println("listening on", so.Addr())

	ed := new(Edit)
	rpc.Register(ed)
	rpc.Accept(so)
}

type Edit int

func (e *Edit) E(req edit.Request, resp *edit.ExitCode) error {
	args := make([]string, len(req.Args))
	for n, arg := range req.Args {
		if strings.Index(arg, "-") == 0 {
			args[n] = arg
		} else {
			path := translatepath(arg)
			args[n] = path
		}
	}
	debug("editing:", args)
	ex, err := E(args...)
	if err != nil {
		debug("unable to start editor:", err)
		return err
	}
	debug("editor exit code:", ex)
	*resp = edit.ExitCode(ex)
	return nil
}

func E(arg ...string) (ex int, err error) {
	e, err := exec.LookPath("E")
	if err != nil {
		return ex, edit.ErrNotFound
	}
	cur, err := os.Executable()
	if err != nil {
		return ex, err
	}
	if edit.IsProcessCircular(cur, e) {
		return ex, edit.ErrCircular
	}
	cmd := exec.Command(e, arg...)
	if err := cmd.Start(); err != nil {
		return ex, edit.ErrCannotExec
	}
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				// that the editor returns a non-zero exit code
				// is not an error in our part
				return status.ExitStatus(), nil
			}
		}
	}
	// ergo, si non non-nulla, nulla est
	return 0, nil
}

func translatepath(remotefilep string) string {
	for _, mountp := range mountps {
		if !mountp.opt("osxfuse") {
			continue
		}
		if strings.HasPrefix(remotefilep, mountp.Remotep) {
			debugf("%s mounted on %s", mountp.Remotep, mountp.Localp)
			localp := strings.Replace(remotefilep, mountp.Remotep, mountp.Localp, 1)
			debugf("translated local file path %s to %s", remotefilep, localp)
			return localp
		}
	}
	debug("path does not have a mountpoint:", remotefilep)
	return remotefilep
}

func debug(v ...interface{}) {
	if *verbose {
		log.Println(v...)
	}
}

func debugf(format string, v ...interface{}) {
	debug(fmt.Sprintf(format, v...))
}
