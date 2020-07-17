edit
====

These are a small set of utilities for editing files from a remote
system over an [ssh(1)] session on your local machine’s text editor.

It currently only supports editing files that are available on the
local system via an [sshfs(1)] mount.

The utilities were written specifically for Plan 9’s [acme(1)]
editor, but could theoretically work with any non-blocking editor
set in `EDITOR_B`.

Install the necessary programs on the respective systems:

	local% go get sny.no/tools/edit/cmd/editd
	remote% go get sny.no/tools/edit/cmd/E sny.no/tools/edit/cmd/B

First start the editd daemon on your local machine:

	% editd &
	./editd: listening on [::]:52670

If you are on macOS you may create a [launch agent] for the editd
daemon in ~/Library/LaunchAgents/editd.plist:

	<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
	<dict>
		<key>Label</key>
		<string>editd</string>
		<key>ProgramArguments</key>
		<array><string>editd<string></array>
		<key>RunAtLoad</key>
		<true />
	</dict>
	</plist>

The agent can be enabled, activated, and deactivated this way:

	% launchctl enable editd
	% launchctl start editd
	% launchctl stop editd

Start an ssh session to the remote machine with a reverse proxy
tunnel for communicating with editd:

	% ssh -R 52670:localhost:52670 <user>@<remote>

You can make ssh do this permanently by adding a section to your
$HOME/.ssh/config file:

	Host <remote>
		RemoteForward 52670 localhost:52670

Use the E or B programs to open files:

	% E hello.txt  # will wait until hello.txt is saved
	% B hello.txt  # will return immediately

_B_ returns immediately which is useful when if you want to open
files for editing and continue using your shell session for other
tasks.  _E_, on the other hand, is blocking and suitable as the
system editor:

	% export EDITOR=$GOBIN/E

Requires OpenSSH 6.7… because of Unix domain socket proxy forwarding.


Known bugs
----------

* Missing support for files not available via a mountpoint on the
  local system.  My current thinking is to remedy this either by
  temporarily mounting and unmounting the remote directory (requires
  [sshfs(1)] and FUSE), or by copying the file to a temporary
  directory on the local end and syncing it back.


[acme(1)]: http://man.cat-v.org/plan_9/1/acme
[ed(1)]: https://manpages.debian.org/buster/ed/ed.1.en.html
[ssh(1)]: https://manpages.debian.org/buster/openssh-client/ssh.1.en.html
[sshfs(1)]: https://manpages.debian.org/buster/sshfs/sshfs.1.en.html

[launch agent]: https://developer.apple.com/library/archive/documentation/MacOSX/Conceptual/BPSystemStartup/Chapters/CreatingLaunchdJobs.html

