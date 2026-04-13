package iotools

import (
	"fmt"
	"encoding/json"
	"os"
)

type Website struct {
	Name		string	`json:"name"`
	URL			string	`json:"url"`
	IsDom		bool	`json:"domestic"`
	StatusOK	bool	`json:"status_ok"`
	Ping 		int64	`json:"ping"`
}

type WebsiteList struct {
	DateTime	string		`json:"datetime"`
	TimeStamp	int64		`json:"ts"`
	Website		[]Website	`json:"website"`
	Elapsed		int64		`json:"elapsed_time"`
	Timeout		int64		`json:"timeout"`
}

func JsonFileToStruct(fname string) (WebsiteList, error) {
	website := WebsiteList{Website: make([]Website, 10, 512)}
	file, err := os.Open(fname)
	if err != nil {
		return website, fmt.Errorf("Error in JsonFileToStruct / Open data file: ", err)
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&website); err != nil {
		return website, fmt.Errorf("Error in JsonFileToStruct / Decoding: ", err)
	}
	return website, nil
}

func StructToJsonFile(fname string, data WebsiteList) error {
	file, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf("Error in StructToJsonFile, file create error", err)
	}
	defer file.Close()
	if err := json.NewEncoder(file).Encode(&data); err != nil {
		return fmt.Errorf("Error in StructToJsonFile, file write error", err)
	}
	return nil
}