package controllers

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"

	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// createSSHSession returns pointers to ssh.Client and ssh.Session, and error if any
func createSSHSession(ip, userName, password, port string) (*ssh.Client, *ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{}
	sshConfig.User = userName
	sshConfig.Auth = []ssh.AuthMethod{}
	if password != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.Password(password))
	}

	if username, err := user.Current(); err == nil {
		userIdRSAFile := username.HomeDir + string(filepath.Separator) + ".ssh" + string(filepath.Separator) + "id_rsa"
		userIdDSAFile := username.HomeDir + string(filepath.Separator) + ".ssh" + string(filepath.Separator) + "id_dsa"

		if _, err := os.Stat(userIdRSAFile); err == nil {
			sshConfig.Auth = append(sshConfig.Auth, PublicKeyFile(userIdRSAFile))
		}
		if _, err := os.Stat(userIdDSAFile); err == nil {
			sshConfig.Auth = append(sshConfig.Auth, PublicKeyFile(userIdDSAFile))
		}
	}

	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", ip+":"+port, sshConfig)
	if err != nil {
		logger.Warn(err.Error())
		return nil, nil, errors.New("cannot dial")
	}

	session, err := client.NewSession()
	if err != nil {
		logger.Warn(err.Error())
		return nil, nil, errors.New("unable to create session")
	}

	return client, session, nil
}
