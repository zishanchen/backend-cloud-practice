package main

import (
	"fmt"
	"os"
)

func basicDefer() {
	fmt.Println("Start")

	defer fmt.Println("Deferred")

	fmt.Println("End")
}

func multupleDefers() {
	defer fmt.Println("First defer")
	defer fmt.Println("Second defer")
	defer fmt.Println("Third defer")

	fmt.Println("Function body")
}

func writeFile() error {
	file, err := os.Create("example.text")

	if err != nil {
		return err
	}

	// Cleanup: close the file when the function finishes
	defer file.Close()

	_, err = file.WriteString("Hello from Go!\n")
	if err != nil {
		return err
	}

	fmt.Println("File written successfully")

	return nil
}

func main() {
	fmt.Println("=== Basic defer ===")
	basicDefer()

	fmt.Println()

	fmt.Println("=== Multiple defers ===")
	multupleDefers()

	fmt.Println()

	fmt.Println("=== Cleanup with defer ===")
	err := writeFile()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Done")
}
