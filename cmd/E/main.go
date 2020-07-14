package main // sny.no/tools/edit/cmd/E

import (
	"flag"
	"fmt"
	"log"
	"os"

	"sny.no/tools/edit"
)

const EX_USAGE  = 64

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
		ex, err := edit.RE(arg...)
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
