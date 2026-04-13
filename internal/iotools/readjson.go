package internal

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
	website := WebsiteList{Website: make([]Website)}
	file, err := os.Open(fnname)
	if err != nil {
		return nil, fmr.Errorf("Error in JsonFileToStr / Open data file: ", err)
	}
	if err := json.NewDecoder(file).Decode(&website); err != nil {
		return nil, fmr.Errorf("Error in JsonFileToStr / Decoding: ", err)
	}
	return website, nil
}