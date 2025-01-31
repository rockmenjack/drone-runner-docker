// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package shell provides functions for converting shell commands
// to posix shell scripts.
package shell

import (
	"bytes"
	"fmt"

	sq "github.com/kballard/go-shellquote"
)

// Script converts a slice of individual shell commands to
// a posix-compliant shell script.
func Script(commands []string) string {
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf)
	fmt.Fprint(buf, optionScript)
	fmt.Fprint(buf, tmateScript)
	fmt.Fprintln(buf)
	for _, command := range commands {
		buf.WriteString(fmt.Sprintf(
			traceScript,
			sq.Join(command),
			command,
		))
	}
	return buf.String()
}

// optionScript is a helper script this is added to the build
// to set shell options, in this case, to exit on error.
const optionScript = `
if [ ! -z "${DRONE_NETRC_FILE}" ]; then
	echo $DRONE_NETRC_FILE > $HOME/.netrc
	chmod 600 $HOME/.netrc
fi

unset DRONE_SCRIPT
unset DRONE_NETRC_MACHINE
unset DRONE_NETRC_USERNAME
unset DRONE_NETRC_PASSWORD
unset DRONE_NETRC_FILE

set -e
`

// traceScript is a helper script that is added to
// the build script to trace a command.
const traceScript = `
echo + %s
%s
`

const tmateScript = `
remote_debug() {
	if [ "$?" -ne "0" ]; then
		/usr/drone/bin/tmate -F
	fi
}

if [ "${DRONE_BUILD_DEBUG}" = "true" ]; then
	if [ ! -z "${DRONE_TMATE_HOST}" ]; then
		echo "set -g tmate-server-host $DRONE_TMATE_HOST" >> $HOME/.tmate.conf
		echo "set -g tmate-server-port $DRONE_TMATE_PORT" >> $HOME/.tmate.conf
		echo "set -g tmate-server-rsa-fingerprint $DRONE_TMATE_FINGERPRINT_RSA" >> $HOME/.tmate.conf
		echo "set -g tmate-server-ed25519-fingerprint $DRONE_TMATE_FINGERPRINT_ED25519" >> $HOME/.tmate.conf

		if [ ! -z "${DRONE_TMATE_AUTHORIZED_KEYS}" ]; then
			echo "$DRONE_TMATE_AUTHORIZED_KEYS" > $HOME/.tmate.authorized_keys
			echo "set -g tmate-authorized-keys \"$HOME/.tmate.authorized_keys\"" >> $HOME/.tmate.conf
		fi
	fi
	trap remote_debug EXIT
fi
`
