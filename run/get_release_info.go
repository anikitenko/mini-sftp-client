package main

import (
	"encoding/json"
	"errors"
	"github.com/bitly/go-simplejson"
	"net/http"
	"runtime"
	"strings"
)

func getReleaseInfo() (*releaseInfo, *staticReleaseInfo, error) {
	resp, err := http.Get(releaseUrl)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	jsonRead, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	tagName, err := jsonRead.Get("tag_name").String()
	if err != nil {
		return nil, nil, err
	}

	info := &releaseInfo{TagName: tagName}

	assets := jsonRead.Get("assets").MustArray()

	currentOs := runtime.GOOS

	for _, v := range assets {
		m, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		assetName, ok := m["name"].(string)
		if !ok {
			continue
		}

		browserDownload, ok := m["browser_download_url"].(string)
		if !ok {
			continue
		}

		assetSizeJSON, ok := m["size"].(json.Number)
		if !ok {
			continue
		}

		assetSize, err := assetSizeJSON.Int64()
		if err != nil {
			continue
		}

		info.Name = assetName
		info.URL = browserDownload

		if strings.Contains(info.Name, "mini-sftp-client-"+currentOs) {
			staticRelease := &staticReleaseInfo{}
			info.Size = int(assetSize)

			for _, v := range assets {
				m, ok := v.(map[string]interface{})
				if !ok {
					continue
				}

				staticName, ok := m["name"].(string)
				if !ok {
					continue
				}
				if strings.Contains(staticName, "update.zip") {
					assetSizeJSON, ok = m["size"].(json.Number)
					if !ok {
						continue
					}

					assetSize, err = assetSizeJSON.Int64()
					if err != nil {
						continue
					}

					staticRelease.Size = int(assetSize)

					browserDownload, ok := m["browser_download_url"].(string)
					if !ok {
						continue
					}

					staticRelease.URL = browserDownload

					return info, staticRelease, nil
				}
			}
		}
	}

	return nil, nil, errors.New("url not found")
}
