package nettools

import (
	"fmt"
	"local/melli/internal/iotools"
	"net/http"
	"time"
	"sync"
	"context"
)

func PingWebsite(ctx context.Context, name string, url string, stdout bool, data chan <- iotools.Website, w *sync.WaitGroup) {
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
	if stdout {
		fmt.Println("Done", url, statusCodeOrError, ping)
	}
	data <- result
}