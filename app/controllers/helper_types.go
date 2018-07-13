package controllers

import (
	"golang.org/x/crypto/ssh"
	"io"
)

type (
	StoredUserPasswordStruct struct {
		User     string
		Password string
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
		Ip       string `json:"ip" example:"127.0.0.1" required:"true"`
		User     string `json:"user" example:"root"`
		Password string `json:"password"`
		Port     string `json:"port" example:"22"`
	}

	GeneralResponse struct {
		Result  bool   `json:"result"`
		Message string `json:"message"`
	}

	GetConnectionsStruct struct {
		Result      bool                  `json:"result"`
		Message     string                `json:"message"`
		Connections []ApiConnectionStruct `json:"connections"`
	}

	GetPathCompletionStruct struct {
		Result  bool     `json:"result"`
		Message string   `json:"message"`
		Items   []string `json:"items"`
	}
)
