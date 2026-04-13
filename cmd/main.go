package main

import (
	"fmt"
	"local/melli/internal/iotools"
	"net/http"
	"time"
)

func PingWebsite(w iotools.Website) (int, int64, error) {
	dt := time.Now().UnixMilli()
	resp, err := http.Get(w.URL)
	dt = time.Now().UnixMilli() - dt
	// req.Header.Add("Content-Type", "application/json")
	// client := &http.Client{}
	// resp, err := client.Do(req)
	if err != nil {
		return -1, dt, fmt.Errorf("Error sending request: %v", err)
	}
	// Ensure the response body is closed.
	defer resp.Body.Close()
	return resp.StatusCode, dt, nil
}


func main() {
	data, _ := iotools.JsonFileToStr("data/d.json")
	fmt.Println(data.Website[0])
	s, t, _ := PingWebsite(data.Website[0])
	fmt.Println(s, t)
}