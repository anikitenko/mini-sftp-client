package controllers

import (
	"io/ioutil"

	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"math"
	"strconv"
)

// CompileJSONResult returns map[string]interface{} from input.
// Function helps to compile result for JSON output
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

// PublicKeyFile returns ssh.AuthMethod which is needed to
// create additional AuthMethod from privte key file
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

// Read implements PassThru struct. Used for displaying copying progress
func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	pt.total += int64(n)

	if err == nil {
		logger.Infof("Read %s for a total of %s", FormatBytes(float64(n)), FormatBytes(float64(pt.total)))
	}

	return n, err
}

// Round returns round of value.
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

// FormatBytes returns formatted string of bytes (1024 bytes -> 1KB)
func FormatBytes(size float64) string {
	if size <= 0 {
		return "Unknown"
	}
	base := math.Log(size) / math.Log(1024)
	var suffixes [5]string
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	if int(math.Floor(base)) > 4 {
		return "Unknown"
	}
	getSuffix := suffixes[int(math.Floor(base))]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix)
}
