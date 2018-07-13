package controllers

import (
	"os/user"
	"strconv"
	"strings"

	"errors"
	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
	"time"
)

// EstablishSSHConnection is a helper function which is used to connect via SSH
// and accepts ssh_ip, ssh_user, ssh_password, ssh_port
func (c App) EstablishSSHConnection() revel.Result {
	sshIPHostname := strings.TrimSpace(c.Params.Get("ssh_ip"))
	sshUser := strings.TrimSpace(c.Params.Get("ssh_user"))
	sshPassword := c.Params.Get("ssh_password")
	sshPort := strings.TrimSpace(c.Params.Get("ssh_port"))

	newConnection := make(map[string][]StoredUserPasswordStruct)

	if errString, err := ConnectSSH(sshIPHostname, sshUser, sshPassword, sshPort); err != nil {
		logger.Warnf("%s: %v", errString, err)
		response := CompileJSONResult(false, errString)
		return c.RenderJSON(response)
	}

	if _, ok := StoredConnection[sshIPHostname]; ok {
		for _, val := range StoredConnection[sshIPHostname] {
			for i, val := range val {
				if i == sshPort {
					for _, val := range val {
						if val.User == sshUser {
							return nil
						}
					}
				}
			}
		}
	}
	newConnection[sshPort] = append(newConnection[sshPort], StoredUserPasswordStruct{User: sshUser, Password: sshPassword})
	StoredConnection[sshIPHostname] = append(StoredConnection[sshIPHostname], newConnection)

	return nil
}

func ConnectSSH(host, username, pass, port string) (string, error) {
	var errorMessage string

	if host == "" {
		errorMessage = "SSH IP is empty"
		return errorMessage, errors.New("ssh ip empty")
	} else if host == MockSSHHostString {
		if !MockSSHServer {
			go createMockSSHServer()
			time.Sleep(time.Second)
			MockSSHServer = true
		}
		host = "127.0.0.1"
	}

	if username == "" {
		localUsername, err := user.Current()
		if err != nil {
			errorMessage = "You didn't specify SSH user and we were not able to determine it from your system"
			return errorMessage, err
		}
		username = localUsername.Username
	}

	if port == "" {
		port = "22"
	} else {
		if _, err := strconv.Atoi(port); err != nil {
			errorMessage = "You specified wrong SSH port"
			return errorMessage, err
		}
	}

	sshSession := createSession(host, username, pass, port)

	SSHclient = sshSession.Client
	SSHsession = sshSession.Session

	if sshSession.ErrorErr != nil {
		return sshSession.ErrorStr, sshSession.ErrorErr
	}

	return "", nil
}

func createSession(host, username, pass, port string) SSHSessionStruct {
	var errorMessage string
	var sshSession SSHSessionStruct
	client, session, err := createSSHSession(host, username, pass, port)

	if err != nil {
		switch err.Error() {
		case "cannot dial":
			errorMessage = "We could not reach '" + host + ":" + port + "' OR login/password is incorrect"
			sshSession.ErrorErr = err
			sshSession.ErrorStr = errorMessage
			return sshSession
		case "unable to create session":
			errorMessage = "We reached '" + host + ":" + port + "' but could not create a test session"
			sshSession.ErrorErr = err
			sshSession.ErrorStr = errorMessage
			return sshSession
		default:
			errorMessage = err.Error()
			sshSession.ErrorErr = err
			sshSession.ErrorStr = errorMessage
			return sshSession
		}
	}

	sshSession.Client = client
	sshSession.Session = session

	return sshSession
}
