edit
====

These are a small set of utilities for editing files from a remote
system over an [ssh(1)] session on your local machine’s text editor.

It should support all editors, but it is written specifically with
Plan 9’s [acme(1)] in mind.

Install the necessary programs on both systems:

	% go get sny.no/tools/edit/...

First start the editd daemon on your local machine:

	% editd & Listening on …

Start an ssh session to the remote machine with a reverse proxy
tunnel for communicating with editd:

	% ssh -R /tmp/edit/ns.am.default:localhost:52670 <user>@<remote>

You can make ssh do this permanently by adding a section to your
$HOME/.ssh/config file:

	Host <remote>
		RemoteForward /tmp/edit/ns.am.default localhost:52670

Use the E or B programs to invoke your favourite editor within this
ssh session:

	% E hello.txt  # will wait until hello.txt is saved % B
	hello.txt  # will return immediately

Requires OpenSSH 6.7… because of Unix domain socket proxy forwarding.

When B or E is invoked on the remote system, editd will translate
the path(s) and call the program set in the `EDITOR` environment
variable.  If `EDITOR` is not defined it will fall back to [ed(1)],
which is the standard editor.


[acme(1)]: http://man.cat-v.org/plan_9/1/acme
[ed(1)]: https://manpages.debian.org/buster/ed/ed.1.en.html
[ssh(1)]: https://manpages.debian.org/buster/openssh-client/ssh.1.en.html
