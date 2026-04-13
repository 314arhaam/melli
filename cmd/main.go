package main

import (
	"fmt"
	"local/melli/internal/iotools"
	// "net/http"
	"time"
	"sync"
	"encoding/json"
	"context"
	"flag"
	"local/melli/internal/nettools"
)

type CLIArgs struct {
	timeOut		*int64
	inputFile	*string
	outputFile	*string
	stdout		bool
}

func (cli *CLIArgs) Init() {
	cli.timeOut = flag.Int64("t", 60, "Timeout in seconds, default 60")
	cli.inputFile = flag.String("f", "", "**REQUIRED** Input filename")
	cli.outputFile = flag.String("o", "", "Output filename. Default: auto generated name with datetime, Pass - for json stdout only.")
	flag.Parse()
	if *cli.inputFile == "" {
		panic("No input file added")
	}
	if *cli.outputFile == "-" {
		cli.stdout = false
	} else {
		cli.stdout = true
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
	out.Start()
	startTime := time.Now()
	for _, wb := range data.Website {
		wg.Add(1)
		go nettools.PingWebsite(ctx, wb.Name, wb.URL, cli.stdout, d, &wg)
	}
	//
	wg.Wait()
	close(d)
	for i := range d {
		out.Website = append(out.Website, i)
	}
	/*out.Elapsed = time.Since(startTime).Milliseconds()
	out.DateTime = startTime.Format(time.DateTime)
	out.TimeStamp = startTime.UnixMilli()*/
	out.Timeout = (*cli.timeOut) * 1000
	//fmt.Println(out)
	out.Tick()
	if !cli.stdout {
		jsonByte, err := json.Marshal(out)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(jsonByte))
		}
	} else {
		var filename string
		if *cli.outputFile != "" {
			filename = *cli.outputFile
		} else {
			// ds, _ := strconv.FormatInt()
			filename = "data/results-" + startTime.Format(time.DateOnly) + "_" + startTime.Format(time.TimeOnly) + ".json"
		}
		err := iotools.StructToJsonFile(filename, out)
		if err != nil {
			fmt.Println(err)
		}
	}
}