package controllers

import (
	"github.com/revel/revel"
	"golang.org/x/crypto/ssh"
	"io"
)

type (
	App struct {
		*revel.Controller
	}

	ApiV1 struct {
		*revel.Controller
	}

	SSHSessionStruct struct {
		Client   *ssh.Client
		Session  *ssh.Session
		ErrorErr error
		ErrorStr string
	}

	FileStructureStruct struct {
		Path      string
		Directory bool
		Symlink   bool
	}

	PassThru struct {
		io.Reader
		total int64
	}

	ApiConnectionStruct struct {
		Ip       string
		User     string
		Password string
		Port     string
	}
)
