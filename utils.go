package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

type pairSlice []pair
type pair struct {
	key int64
	val string
}

func (p pairSlice) Len() int {
	return len(p)
}

func (p pairSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p pairSlice) Less(i, j int) bool {
	return p[i].key < p[j].key
}

func parseFloatOrFail(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		logfile.Fatalf("could not parse %s\n", str)
	}
	return val
}

func logger(file string) *log.Logger {
	f, err := os.Create(file)
	if err != nil {
		logfile.Fatalf("failed to create log file %s: %s\n", file, err)
	}
	return log.New(f, "", log.LstdFlags)
}

func write(w *bufio.Writer, str string) {
	_, err := w.WriteString(str)
	if err != nil {
		logfile.Fatalf("trouble writing string: %s", err)
	}
}
