package main

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
)

type mountpoint struct {
	Remote  string
	Remotep string
	Localp  string
	Opts    []string
}

func (p *mountpoint) opt(name string) bool {
	for _, opt := range p.Opts {
		if opt == name {
			return true
		}
	}
	return false
}

func mountpoints() ([]mountpoint, error) {
	cmd := exec.Command("mount")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(out)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	ps := make([]mountpoint, 0)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		p := parsemount(line)
		for _, opt := range p.Opts {
			if opt == "osxfuse" {
				ps = append(ps, p)
			}
		}
	}
	return ps, nil
}

// andretol@job.uio.no:/uio/kant/usit-uait-u1/andretol on /Volumes/uio (osxfuse, nodev, nosuid, synchronous, mounted by am)

// "<remote>:<remotep> on <localp> (<opt>[, <opt>])"
func parsemount(line string) mountpoint {
	// TODO: replace with a parser
	if strings.Index(line, "osxfuse") < 0 {
		return mountpoint{}
	}
	fields := strings.Fields(line)
	f1 := strings.Split(fields[0], ":")
	return mountpoint{
		Remote:  f1[0],
		Remotep: f1[1],
		Localp:  fields[2],
		Opts:    []string{"osxfuse"},
	}
}
