package main

import (
	"bufio"
	"fmt"
	"os"
)

func writeFile(filename string, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func appendFile(filename string, content string) error {
	file, err := os.OpenFile(
		filename,
		os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("open file for append: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("append file: %w", err)
	}

	return nil
}

func readFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	fmt.Println("File content:")

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	return nil
}

func main() {
	filename := "example.txt"

	err := writeFile(
		filename,
		"Hello from Go!\n",
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = appendFile(
		filename,
		"This line was appended.\n",
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = readFile(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("File operations completed successfully.")
}
