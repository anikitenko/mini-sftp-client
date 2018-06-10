package controllers

import (
	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"math/rand"
	"strconv"
	"time"
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
	PinCode = strconv.Itoa(rand.Intn(9999))
	logger.Warnf("Your pin code: %s. You will need this in order to access client not from your local machine!", PinCode)
}

func checkPinCode(c *revel.Controller) revel.Result {
	r := c.Request

	if r.Method != "POST" || c.ClientIP != "127.0.0.1" {
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
