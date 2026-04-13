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

func PingWebsite(ctx context.Context, name string, url string, data chan <- iotools.Website, w *sync.WaitGroup) {
	defer w.Done()
	startTime := time.Now()
	status := true
	statusCodeOrError := "200 OK"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		status = false
		statusCodeOrError = err.Error()
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	ping := time.Since(startTime).Milliseconds()
	if err != nil || resp.StatusCode != 200 {
		status = false
		statusCodeOrError = err.Error()
	} else {
		defer resp.Body.Close()
	}
	result := iotools.Website{
		Name: name,
		URL: url, 
		Ping: ping, 
		StatusOK: status,
	}
	fmt.Println("Done", url, statusCodeOrError, ping)
	data <- result
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
	d := make(chan iotools.Website, len(data.Website))
	// Data from channels aggregated
	out := iotools.WebsiteList{}
	//
	startTime := time.Now()
	for _, wb := range data.Website {
		wg.Add(1)
		go PingWebsite(ctx, wb.Name, wb.URL, d, &wg)
	}
	//
	wg.Wait()
	close(d)
	for i := range d {
		out.Website = append(out.Website, i)
	}
	out.Elapsed = time.Since(startTime).Milliseconds()
	out.DateTime = startTime.Format(time.DateTime)
	out.TimeStamp = startTime.UnixMilli()
	//fmt.Println(out)
	err := iotools.StructToJsonFile(*cli.outputFile, out)
	if err != nil {
		fmt.Println(err)
	}
}