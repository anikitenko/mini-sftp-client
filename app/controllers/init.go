package controllers

import (
	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"time"
	"io"
	"crypto/rand"
)

// Defines global variables
var (
	SSHclient            *ssh.Client
	SSHsession           *ssh.Session
	MockSSHServer        = false
	MockSSHHostString    = "sftp-mock-test"
	MockSSHUser          = "test"
	MockSSHPass          = "test"
	PinCode              string
	TimeToWaitInvalidPin time.Duration
)

func GeneratePinCode() {
	table := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, 4)
	_, err := io.ReadAtLeast(rand.Reader, b, 4)
	if err != nil {
		logger.Fatalf("Problem with generating pin code: %v", err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	PinCode = string(b)
	logger.Warnf("Your pin code: %s. You will need this in order to access client not from your local machine!", PinCode)
}

func checkPinCode(c *revel.Controller) revel.Result {
	r := c.Request

	if r.Method != "POST" || c.ClientIP == "127.0.0.1" {
		return nil
	}

	userPinCode := c.Session["pin_code"]
	if PinCode != userPinCode {
		return c.Forbidden("You are not permitted to make this request")
	}

	return nil
}

func init() {
	revel.InterceptFunc(checkPinCode, revel.BEFORE, &App{})
}
