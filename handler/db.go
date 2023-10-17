package handler

import (
	"encoding/json"
	"os"
)

type Url struct {
	Id     string `json:"id"`
	SSL    bool   `json:"ssl"`
	Opaque string `json:"opaque"` // encoded opaque data
	Host   string `json:"host"`   // host or host:port
	Path   string `json:"path"`   // path (relative paths may omit leading slash)
	Query  string `json:"query"`  // encoded query values, without '?'
}

var db_file = "./data/db.json"

func GetUrls() (urls map[string]Url) {
	urls = map[string]Url{}
	data, err := os.ReadFile(db_file)
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
	err = os.WriteFile(db_file, data, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
