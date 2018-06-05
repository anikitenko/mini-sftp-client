package controllers

import (
	"io/ioutil"

	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func CompileJSONResult(result bool, message string, otherData ...map[string]interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	data["result"] = result
	data["message"] = message

	for _, dataOption := range otherData {
		for key, val := range dataOption {
			data[key] = val
		}
	}

	return data
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Warn(err.Error())
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		logger.Warn(err.Error())
		return nil
	}
	return ssh.PublicKeys(key)
}
