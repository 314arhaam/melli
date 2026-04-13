package main

import (
	"fmt"
	"local/melli/internal/iotools"
	"net/http"
	"time"
	"sync"
	"encoding/json"
	"os"
)

func PingWebsite(website <- chan string, data chan <- iotools.Website, w *sync.WaitGroup) {
	defer w.Done()
	t0 := time.Now().UnixMilli()
	for wb := range website {
		resp, err := http.Get(wb)
		t1 := time.Now().UnixMilli()
		ping := int(t1-t0)
		status := true
		if err != nil || resp.StatusCode != 200 {
			status = false
		}
		data <- iotools.Website{URL: wb, Ping: ping, StatusOK: status,}
	}
}


func main() {
	data, _ := iotools.JsonFileToStr("data/d.json")
	fmt.Println(data.Website)
	w := make(chan string, len(data.Website))
	d := make(chan iotools.Website, len(data.Website))
	var wg sync.WaitGroup
	wg.Add(1)
	dt := time.Now().UnixMilli()
	go PingWebsite(w, d, &wg)
	/*
	go func(d <- chan iotools.Website){
		for data := range d {
			fmt.Println(data)
		}
	}(d)
	*/
	for i := 0; i < len(data.Website); i++ {
		w <- data.Website[i].URL
	}
	close(w)
	wg.Wait()
	close(d)
	out := iotools.WebsiteList{}
	for i := range d {
		fmt.Println("Done", i)
		out.Website = append(out.Website, i)
	}
	dt = time.Now().UnixMilli() - dt
	out.Elapsed = int(dt)
	fmt.Println(out)
	file, err := os.Create("data/results.json")
	if err != nil {
		fmt.Errorf("Error in main, file creation error", err)
	}
	defer file.Close()
	if err := json.NewEncoder(file).Encode(&out); err != nil {
		fmt.Println(err)
	}
}