package main

import (
	"encoding/json"
	"errors"
	"github.com/bitly/go-simplejson"
	"net/http"
	"runtime"
	"strings"
)

func getReleaseInfo() (*releaseInfo, error) {
	resp, err := http.Get(releaseUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	jsonRead, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	tagName, err := jsonRead.Get("tag_name").String()
	if err != nil {
		return nil, err
	}

	info := &releaseInfo{TagName: tagName}

	assets := jsonRead.Get("assets").MustArray()

	currentOs := runtime.GOOS

	for _, v := range assets {
		m, ok := v.(map[string]interface{})
		if !ok {
			_ = errors.New("invalid data")
		}
		assetName := m["name"].(string)
		browserDownload := m["browser_download_url"].(string)
		assetSize, err := m["size"].(json.Number).Int64()
		if err != nil {
			return nil, err
		}
		info.Name = assetName
		info.URL = browserDownload

		if strings.Contains(info.Name, "mini-sftp-client-"+currentOs) {
			info.Size = int(assetSize)
			return info, nil
		}
	}

	return nil, errors.New("url not found for os")
}
