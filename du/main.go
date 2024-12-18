package main

import (
	"flag"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

var dir = flag.String("dir", ".", "path to folder")
var units = flag.String("units", "b", "result units: b, kb, gb.")
var output = flag.Bool("list-files", false, "show files list")

func main() {
	flag.Parse()

	sizeInBytes := totalSize(*dir)

	fmt.Printf("Total size: %s\n", formatSize(sizeInBytes, *units))
}

func totalSize(path string) int64 {
	var total int64
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		info, err := d.Info()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if *output {
			fmt.Printf("%s:%d\n", d.Name(), info.Size())
		}
		total += info.Size()
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return total
}

func formatSize(sizeInBytes int64, targetUnits string) string {
	switch strings.ToLower(targetUnits) {
	case "kb":
		return fmt.Sprintf("%.2f Kb", float64(sizeInBytes)/1024)
	case "mb":
		return fmt.Sprintf("%.2f Mb", float64(sizeInBytes)/1024/1024)
	case "gb":
		return fmt.Sprintf("%.2f Gb", float64(sizeInBytes)/1024/1024/1024)
	case "b":
		return fmt.Sprintf("%d bytes", sizeInBytes)
	}
	return fmt.Sprintf("%d bytes", sizeInBytes)
}
