package iotools

import (
	"fmt"
	"encoding/json"
	"os"
)

type Website struct {
	Name	string	`json:"name"`
	URL		string	`json:"url"`
	IsDom	bool	`json:"domestic"`
}

type WebsiteList struct {
	Website	[]Website	`json:"website"`
}

func JsonFileToStr(fname string) (WebsiteList, error) {
	website := WebsiteList{Website: make([]Website, 10, 512)}
	file, err := os.Open(fname)
	if err != nil {
		return website, fmt.Errorf("Error in JsonFileToStr / Open data file: ", err)
	}
	if err := json.NewDecoder(file).Decode(&website); err != nil {
		return website, fmt.Errorf("Error in JsonFileToStr / Decoding: ", err)
	}
	return website, nil
}