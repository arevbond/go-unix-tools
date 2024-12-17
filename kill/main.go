package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
)

var pid = flag.Int("pid", -1, "process identifier")
var signal = flag.Int("signal", 2, "select a signal to send")

var signals = [...]string{
	1:  "hangup",
	2:  "interrupt",
	3:  "quit",
	4:  "illegal instruction",
	5:  "trace/breakpoint trap",
	6:  "aborted",
	7:  "bus error",
	8:  "floating point exception",
	9:  "killed",
	10: "user defined signal 1",
	11: "segmentation fault",
	12: "user defined signal 2",
	13: "broken pipe",
	14: "alarm clock",
	15: "terminated",
	16: "stack fault",
	17: "child exited",
	18: "continued",
	19: "stopped (signal)",
	20: "stopped",
	21: "stopped (tty input)",
	22: "stopped (tty output)",
	23: "urgent I/O condition",
	24: "CPU time limit exceeded",
	25: "file size limit exceeded",
	26: "virtual timer expired",
	27: "profiling timer expired",
	28: "window changed",
	29: "I/O possible",
	30: "power failure",
	31: "bad system call",
}

func main() {
	flag.Parse()

	process, err := os.FindProcess(*pid)
	if err != nil {
		log.Fatal(err)
	}

	err = process.Signal(syscall.Signal(0))
	if err != nil {
		log.Fatal(err)
	}

	err = process.Signal(syscall.Signal(*signal))
	if err != nil {
		log.Fatal(err)
	}

	if *signal < len(signals) {
		fmt.Printf("success send signal %d (%s) to a proccess with pid %d\n",
			*signal, signals[*signal], *pid)
	} else {
		fmt.Printf("success send signal %d to a proccess with pid %d\n",
			*signal, *pid)
	}
}
