package controllers

import "golang.org/x/crypto/ssh"

// Defines global variables
var (
	SSHclient         *ssh.Client
	SSHsession        *ssh.Session
	MockSSHServer     = false
	MockSSHHostString = "sftp-mock-test"
	MockSSHUser       = "test"
	MockSSHPass       = "test"
)
