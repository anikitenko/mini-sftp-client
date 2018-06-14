package controllers

import (
	"github.com/gliderlabs/ssh"
	logger "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

func createMockSSHServer() {
	s := &ssh.Server{
		Addr: ":2222",
		PasswordHandler: func(ctx ssh.Context, pass string) bool {
			if ctx.User() == MockSSHUser && pass == MockSSHPass {
				return true
			}
			return false
		},
	}
	s.Handle(func(s ssh.Session) {
		cmd := exec.Command("bash", "-c", strings.Join(s.Command(), " "))
		out, err := cmd.CombinedOutput()
		if err != nil {
			logger.Warnf("Problem with executing command: %v", err)
			return
		}
		if _, err := s.Write(out); err != nil {
			logger.Warnf("Problem with writing output: %v", err)
			return
		}
	})
	if err := s.ListenAndServe(); err != nil {
		MockSSHServer = false
		logger.Warnf("Problem with running mock SSH server: %v", err)
	}
}
