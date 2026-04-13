package main

import (
	"fmt"
	"local/melli/internal/iotools"
	"net/http"
	"time"
	"sync"
	// "encoding/json"
	"context"
	"flag"
)

func PingWebsite(ctx context.Context, website <- chan string, data chan <- iotools.Website, w *sync.WaitGroup) {
	defer w.Done()
	t0 := time.Now().UnixMilli()
	for url := range website {
		status := true
		statusCodeOrError := "200 OK"
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			status = false
			statusCodeOrError = err.Error()
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		t1 := time.Now().UnixMilli()
		if err != nil || resp.StatusCode != 200 {
			status = false
			statusCodeOrError = err.Error()
		} else {
			defer resp.Body.Close()
		}
		ping := int(t1-t0)
		result := iotools.Website{URL: url, Ping: ping, StatusOK: status,}
		fmt.Println("Done", url, statusCodeOrError, ping)
		data <- result
	}
}

type CLIArgs struct {
	timeOut		*int64
	inputFile	*string
	outputFile	*string
}

func (cli *CLIArgs) Init() {
	cli.timeOut = flag.Int64("t", 60, "Timeout in seconds, default 60")
	cli.inputFile = flag.String("f", "", "Input filename")
	cli.outputFile = flag.String("o", "data/output.json", "Output filename")
	flag.Parse()
	if *cli.inputFile == "" {
		panic("No input file added")
	}
}

func main() {
	// Parse CLI Args
	cli := CLIArgs{}
	cli.Init()
	// Input data
	data, _ := iotools.JsonFileToStruct(*cli.inputFile)
	// Context
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(*cli.timeOut) * time.Second)
	// WaitGroup
	var wg sync.WaitGroup
	// Channels
	w := make(chan string, len(data.Website))
	d := make(chan iotools.Website, len(data.Website))
	// Data from channels aggregated
	out := iotools.WebsiteList{}
	//
	wg.Add(1)
	t0 := time.Now().UnixMilli()
	go PingWebsite(ctx, w, d, &wg)
	for i := 0; i < len(data.Website); i++ {
		w <- data.Website[i].URL
	}
	close(w)
	wg.Wait()
	close(d)
	for i := range d {
		out.Website = append(out.Website, i)
	}
	t1 := time.Now().UnixMilli()
	out.Elapsed = int(t1-t0)
	fmt.Println(out)
	err := iotools.StructToJsonFile(*cli.outputFile, out)
	if err != nil {
		fmt.Println(err)
	}
}