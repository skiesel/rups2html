package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const (
	HistoryFile      = "rups.history"
	MaxHistoryPoints = 50
	DateDelimiter    = "DATEDELIM"
)

type RupsHistory struct {
	machines   map[string]int64
	timestamps []string
	usages     [][]string
	pointCount int64
}

func newRupsHistory() *RupsHistory {
	history := new(RupsHistory)
	history.machines = make(map[string]int64, 0)
	history.timestamps = make([]string, 0)
	history.usages = make([][]string, 0)
	return history
}

func readRupsHistory() (history *RupsHistory) {
	history = newRupsHistory()
	historyFile, err := os.Open(HistoryFile)
	if err != nil {
		return
	}

	defer historyFile.Close()

	lineNumber := 0
	scanner := bufio.NewScanner(historyFile)
	for scanner.Scan() {
		if lineNumber == 0 {
			historyLine := strings.Split(scanner.Text(), DateDelimiter)
			history.loadTimestamps(historyLine)
		} else {
			historyLine := strings.Split(scanner.Text(), " ")
			history.loadMachineHistory(historyLine)
			history.pointCount++
		}
		lineNumber++
	}

	history.checkAndFixHistorySize()
	return
}

func (rh *RupsHistory) saveRupsHistory() {
	outfile, err := os.Create(HistoryFile)
	if err != nil {
		logfile.Fatalf("failed to create %s (%s)", HistoryFile, err)
	}
	historyWriter := bufio.NewWriter(outfile)

	defer func() {
		historyWriter.Flush()
		outfile.Close()
	}()

	timestamps := strings.Join(rh.timestamps, DateDelimiter) + "\n"
	write(historyWriter, timestamps)

	ps := make(pairSlice, 0, len(rh.machines))
	for id, index := range rh.machines {
		ps = append(ps, pair{key: index, val: id})
	}
	sort.Sort(ps)

	for dsIndex, p := range ps {
		usage := strings.Join(rh.usages[dsIndex], " ")
		machineLine := fmt.Sprintf("%s %s\n", p.val, usage)
		write(historyWriter, machineLine)
	}
}

func (rh *RupsHistory) addMachine(id string) {
	_, found := rh.machines[id]
	if !found {
		rh.machines[id] = int64(len(rh.usages))
		rh.usages = append(rh.usages, make([]string, 0))
	} else {
		logfile.Fatalf("Already found: %s\n", id)
	}
}

func (rh *RupsHistory) loadTimestamps(timestampLine []string) {
	for _, val := range timestampLine {
		rh.timestamps = append(rh.timestamps, val)
	}
}

func (rh *RupsHistory) loadMachineHistory(historyLine []string) {
	var machineId string
	for i, val := range historyLine {
		if i == 0 {
			machineId = val
		} else {
			rh.addHistoryPoint(machineId, val)
		}
	}
}

func (rh *RupsHistory) addHistoryPoint(id, val string) {
	index, found := rh.machines[id]
	if !found {
		rh.addMachine(id)
		index, _ = rh.machines[id]
	}
	rh.usages[index] = append(rh.usages[index], val)
}

func (rh *RupsHistory) addCurrentRups() {
	c := exec.Command("/bin/sh", "-c", "rups")
	output, err := c.CombinedOutput()
	if err != nil {
		logfile.Fatalf("failed output=[%s]: %s\n", output, err)
	}
	buf := bytes.NewBuffer(output)

	lines := strings.Split(buf.String(), "\n")
	for i, line := range lines {
		if i == 0 { //dateline
			rh.addNewTimestamp(line)
		} else { //cpu lines
		       	if line != "" {
			   rh.addNewMachineHistoryPoint(line)
			}
		}
	}
	rh.checkAndFixHistorySize()
}

func (rh *RupsHistory) checkAndFixHistorySize() {
	max := 0
	for _, usage := range rh.usages {
		if max < len(usage) {
			max = len(usage)

			if max >= MaxHistoryPoints {
				max = MaxHistoryPoints
				break
			}
		}
	}

	rh.pointCount = int64(max)

	for i, usage := range rh.usages {
		if len(usage) < max {
			for len(usage) < max {
				usage = append(usage, "-1")
			}
			rh.usages[i] = usage
		}
	}

	if len(rh.usages[0]) > MaxHistoryPoints {
		start := len(rh.usages[0]) - MaxHistoryPoints
		for i, usage := range rh.usages {
			rh.usages[i] = usage[start:]
		}
	}

	if len(rh.timestamps) >= MaxHistoryPoints {
		start := len(rh.timestamps) - MaxHistoryPoints
		rh.timestamps = rh.timestamps[start:]
	}
}

func (rh *RupsHistory) addNewTimestamp(timestamp string) {
	rh.timestamps = append(rh.timestamps, timestamp)
}

func (rh *RupsHistory) addNewMachineHistoryPoint(machineLine string) {
	//ai1.cs.unh.edu        17:01 up   2 days,    8:04,    0 user, load 0.00 0.00 0.00
	//    0                   1    2     3          4      5   6     7    8    9   10
	tokens := strings.Split(machineLine, " ")
	machine := tokens[0]
	machineTokens := strings.Split(machine, ".")
	machineId := machineTokens[0]

	count := 0
	avgL := 0.
	for _, token := range tokens {
		switch count {
		case 8:
		case 9:
		case 10:
			val, err := strconv.ParseFloat(token, 64)
			if err != nil {
				// This machine will be missing a data point, but
				// we'll resolve this in checkAndFixHistorySize
				return
			}
			avgL += val
		}
		if token != "" {
			count++
		}
	}

	avgL /= 3

	floatString := strconv.FormatFloat(avgL, 'f', -1, 64)

	rh.addHistoryPoint(machineId, floatString)
}

func (rh *RupsHistory) dumpText() {
	for _, timestamp := range rh.timestamps {
		fmt.Printf("\t%s", timestamp)
	}

	ps := make(pairSlice, 0, len(rh.machines))
	for id, index := range rh.machines {
		ps = append(ps, pair{key: index, val: id})
	}
	fmt.Printf("\n")

	sort.Sort(ps)
	for _, p := range ps {
		fmt.Printf("%d) %s", p.key, p.val)
		for _, val := range rh.usages[p.key] {
			fmt.Printf("\t%f", val)
		}
		fmt.Printf("\n")
	}
}
