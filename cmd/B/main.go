package main // sny.no/tools/edit/cmd/B

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const EX_USAGE = 64

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s file...\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(EX_USAGE)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("%s: ", os.Args[0]))
	flag.Usage = usage
	flag.Parse()

	log.Println("B!")
}
