package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	resp, err := http.Get("https://example.com")
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
	fmt.Println("Status Code:", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ready body error:", err)
		return
	}

	fmt.Println("Body length:", len(body))
	fmt.Println("Body preview:")

	previewLength := 200
	if len(body) < previewLength {
		previewLength = len(body)
	}

	fmt.Println(string(body[:previewLength]))
}
