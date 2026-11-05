package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func fetch(url string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set(
		"User-Agent",
		"study-go-http-client/1.0",
	)

	req.Header.Set(
		"Accept",
		"text/html,application/json",
	)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println("\n=== Response ===")
	fmt.Println("Status:", resp.Status)
	fmt.Println("Status Code:", resp.StatusCode)
	fmt.Println(
		"Content-Type:",
		resp.Header.Get("Content-Type"),
	)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	fmt.Println("Body Size:", len(body), "bytes")

	if resp.StatusCode < 200 ||
		resp.StatusCode >= 300 {

		return fmt.Errorf(
			"unexpected HTTP status: %s",
			resp.Status,
		)
	}

	fmt.Println("\n=== Body Preview ===")

	const maxPreview = 500

	previewLength := len(body)

	if previewLength > maxPreview {
		previewLength = maxPreview
	}

	fmt.Println(string(body[:previewLength]))

	if len(body) > maxPreview {
		fmt.Println("\n... response truncated")
	}

	return nil
}

func main() {
	url := "https://example.com"

	if len(os.Args) > 1 {
		url = os.Args[1]
	}

	fmt.Println("GET:", url)

	start := time.Now()

	err := fetch(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(
		"\nCompleted in:",
		time.Since(start),
	)
}
