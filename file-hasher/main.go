package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var dir = flag.String("dir", ".", "target directory")
var workers = flag.Int("workers", 1, "amount parallel goroutines")
var algorithm = flag.String("algorithm", "sha256", "crypto algorithm")
var recursive = flag.Bool("recursive", false, "walk all files")

type Task struct {
	filaPath  string
	algorithm string
}

type Result struct {
	filaPath string
	hash     string
}

func main() {
	start := time.Now()

	flag.Parse()

	fmt.Printf("Started hashing files in %s using %d workers\n", *dir, *workers)

	results := make(chan Result)

	tasks := make(chan Task)

	var totalFiles int
	if *recursive {
		go func() {
			err := filepath.WalkDir(*dir, func(path string, d fs.DirEntry, err error) error {
				fmt.Println(path + d.Name())
				task := Task{filaPath: path + d.Name(), algorithm: *algorithm}
				tasks <- task
				totalFiles++
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}
			close(tasks)
		}()

	} else {
		files, err := os.ReadDir(*dir)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				task := Task{filaPath: *dir + "/" + file.Name(), algorithm: *algorithm}
				tasks <- task
				totalFiles++
			}
			close(tasks)
		}()

	}

	wg := &sync.WaitGroup{}
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go generateHash(tasks, results, wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	cnt := 1
	for result := range results {
		fmt.Printf("[%d/%d] %s: %s\n", cnt, totalFiles, result.filaPath, result.hash)
		cnt++
	}
	fmt.Printf("Complete! Processed %d files in %f seconds\n", totalFiles, time.Since(start).Seconds())
}

func generateHash(tasks chan Task, results chan Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		f, err := os.Open(task.filaPath)
		if err != nil {
			log.Printf("[ERROR] can't open file: %s: %s", task.filaPath, err)
		}

		var h hash.Hash
		switch strings.ToLower(task.algorithm) {
		case "sha1":
			h = sha1.New()
		case "md5":
			h = md5.New()
		default:
			h = sha256.New()
		}
		if _, err := io.Copy(h, f); err != nil {
			log.Println("[ERROR]", err)
		}
		hexHash := hex.EncodeToString(h.Sum(nil))
		result := Result{
			filaPath: task.filaPath,
			hash:     hexHash,
		}
		results <- result
		err = f.Close()
		if err != nil {
			log.Println("[ERROR] can't close file", err)
		}
	}
}
