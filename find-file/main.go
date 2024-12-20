package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sync"
	"time"
)

var dir = flag.String("dir", ".", "target directory")
var name = flag.String("name", "", "file name to search")
var workers = flag.Int("workers", 1, "amout of goroutines")

var wg sync.WaitGroup
var walWg sync.WaitGroup

func main() {
	start := time.Now()
	flag.Parse()

	fileNames := make(chan string)

	go func() {
		err := filepath.WalkDir(*dir, func(path string, d fs.DirEntry, err error) error {
			fileNames <- path
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		close(fileNames)
	}()

	results := make(chan string)

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go matcher(fileNames, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("Found:")
	for result := range results {
		fmt.Println(result)
	}
	fmt.Printf("\nComplete! Time spent: %s\n", time.Since(start).String())
}

func matcher(fileNames chan string, results chan string) {
	defer wg.Done()
	for fileName := range fileNames {
		if filepath.Base(fileName) == *name {
			results <- fileName
		}
	}
}
