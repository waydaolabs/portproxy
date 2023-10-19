package handler

import (
	"encoding/json"
	"os"

	"github.com/waydaolabs/portproxy/config"
)

type Url struct {
	Id     string `json:"id"`
	SSL    bool   `json:"ssl"`
	Opaque string `json:"opaque"` // encoded opaque data
	Host   string `json:"host"`   // host or host:port
	Path   string `json:"path"`   // path (relative paths may omit leading slash)
	Query  string `json:"query"`  // encoded query values, without '?'
}

func GetUrls() (urls map[string]Url) {
	urls = map[string]Url{}
	data, err := os.ReadFile(config.DB_FILE)
	if err == nil {
		json.Unmarshal(data, &urls)
	}
	return
}

func SetUrls(u Url) {
	urls := GetUrls()
	urls[u.Id] = u
	data, err := json.Marshal(urls)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(config.DB_FILE, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
