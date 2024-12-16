package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"text/tabwriter"
	"unicode"
)

var sortedOutput = flag.Bool("sorted", false, "sort process")

func main() {
	flag.Parse()
	const procPath = "/proc"

	err := os.Chdir(procPath)
	if err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	procsFiles := make([]os.DirEntry, 0)

	for _, file := range files {
		if isProcess(file.Name()) {
			procsFiles = append(procsFiles, file)
		}
	}

	if *sortedOutput {
		sort.Slice(procsFiles, func(i, j int) bool {
			num1, err := strconv.Atoi(procsFiles[i].Name())
			if err != nil {
				log.Fatal(err)
			}
			num2, err := strconv.Atoi(procsFiles[j].Name())
			if err != nil {
				log.Fatal(err)
			}
			return num1 < num2
		})
	}

	type process struct {
		name     string
		pid      int
		username string
	}

	procs := make([]*process, 0, len(procsFiles))

	for _, file := range procsFiles {
		pid, err := strconv.Atoi(file.Name())
		if err != nil {
			log.Fatal(err)
		}

		proc := &process{pid: pid}

		data, err := os.ReadFile(fmt.Sprintf("%s/%s", file.Name(), "comm"))
		if err != nil {
			log.Fatal(err)
		}
		proc.name = string(data)
		procs = append(procs, proc)
	}

	w := tabwriter.NewWriter(os.Stdout, 20, 30, 1, '\t', 0)
	fmt.Fprintf(w, "%s\t%s\n", "PID", "COMMAND")
	for _, p := range procs {
		fmt.Fprintf(w, "%d\t%s\n", p.pid, p.name)
	}

}

func isProcess(name string) bool {
	for _, ch := range name {
		if !unicode.IsDigit(ch) {
			return false
		}
	}
	return true
}
