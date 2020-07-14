package edit // sny.no/tools/edit

type Request struct {
	Args    []string
	Environ []string
}

// Exit code from editor, relayed from remote host or from local process.
type ExitCode int
