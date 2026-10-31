package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type LogStats struct {
	Total  int
	Levels map[string]int
	Errors map[string]int
}

func parseLogFile(filename string) (LogStats, error) {
	file, err := os.Open(filename)
	if err != nil {
		return LogStats{}, fmt.Errorf("open log file: %w", err)
	}
	defer file.Close()

	stats := LogStats{
		Levels: map[string]int{
			"INFO":  0,
			"WARN":  0,
			"ERROR": 0,
		},
		Errors: make(map[string]int),
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		stats.Total++

		switch {
		case strings.HasPrefix(line, "INFO "):
			stats.Levels["INFO"]++

		case strings.HasPrefix(line, "WARN "):
			stats.Levels["WARN"]++

		case strings.HasPrefix(line, "ERROR "):
			stats.Levels["ERROR"]++

			message := strings.TrimPrefix(line, "ERROR ")
			stats.Errors[message]++
		}
	}

	if err := scanner.Err(); err != nil {
		return LogStats{}, fmt.Errorf("scan log file: %w", err)
	}

	return stats, nil
}

func printSummary(stats LogStats) {
	fmt.Println("\n=== Log Summary ===")
	fmt.Println("Total:", stats.Total)
	fmt.Println("INFO:", stats.Levels["INFO"])
	fmt.Println("WARN:", stats.Levels["WARN"])
	fmt.Println("ERROR:", stats.Levels["ERROR"])
}

func printTopErrors(errors map[string]int) {
	if len(errors) == 0 {
		fmt.Println("\nNo errors found.")
		return
	}

	type ErrorCount struct {
		Message string
		Count   int
	}

	errorList := make([]ErrorCount, 0, len(errors))

	for message, count := range errors {
		errorList = append(errorList, ErrorCount{
			Message: message,
			Count:   count,
		})
	}

	sort.Slice(errorList, func(i, j int) bool {
		return errorList[i].Count > errorList[j].Count
	})

	fmt.Println("\n=== Top Errors ===")

	for i, item := range errorList {
		fmt.Printf("%d. %s (%d)\n", i+1, item.Message, item.Count)
	}
}

func main() {
	filename := "cli/log-parser/sample.log"

	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	fmt.Println("Parsing:", filename)

	stats, err := parseLogFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	printSummary(stats)
	printTopErrors(stats.Errors)
}
