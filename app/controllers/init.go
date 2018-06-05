package controllers

import "golang.org/x/crypto/ssh"

var (
	SSHclient  *ssh.Client
	SSHsession *ssh.Session
)
