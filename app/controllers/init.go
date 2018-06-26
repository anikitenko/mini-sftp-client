package controllers

import (
	"crypto/rand"
	"github.com/revel/revel"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"time"
	"net/http"
	"github.com/swaggo/http-swagger"
	_ "mini-sftp-client/app/controllers/docs"
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
	ApiConnections       = make(map[string]ApiConnectionStruct)
)

func installHandlers() {
	revel.AddInitEventHandler(func(event int, _ interface{}) (r int) {
		if event==revel.ENGINE_STARTED {
			var (
				serveMux     = http.NewServeMux()
				revelHandler = revel.CurrentEngine.(*revel.GoHttpServer).Server.Handler
			)

			serveMux.Handle("/",     revelHandler)
			serveMux.Handle("/api/v1/index.html", httpSwagger.WrapHandler)
			serveMux.Handle("/api/v1/doc.json", httpSwagger.WrapHandler)
			revel.CurrentEngine.(*revel.GoHttpServer).Server.Handler = serveMux
		}
		return
	})
}

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
	testParam := c.Params.Get("for_testing")
	r := c.Request

	if c.ClientIP == "127.0.0.1" && testParam != "true" {
		if r.Method != "POST" {
			return nil
		}
	}

	if r.Method != "POST" || (c.ClientIP == "127.0.0.1" && testParam != "true") {
		return nil
	}

	userPinCode := c.Session["pin_code"]
	if PinCode != userPinCode {
		return c.Forbidden("You are not permitted to make this request")
	}

	return nil
}

func checkApiPinCode(c *revel.Controller) revel.Result {
	testParam := c.Params.Get("for_testing")
	r := c.Request

	if c.ClientIP == "127.0.0.1" && testParam != "true" {
		return nil
	}

	userPinCode := r.Header.Get("Pin-Code")
	if PinCode != userPinCode {
		return c.Forbidden("You are not permitted to make this request")
	}

	return nil
}

func init() {
	revel.OnAppStart(installHandlers)
	revel.OnAppStart(GeneratePinCode)
	revel.InterceptFunc(checkPinCode, revel.BEFORE, &App{})
	revel.InterceptFunc(checkApiPinCode, revel.BEFORE, &ApiV1{})
}
