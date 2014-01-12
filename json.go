package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	JSONFile = "html/js/rups.json"
)

func (rh *RupsHistory) createJSONFile() {
	outfile, err := os.Create(JSONFile)
	if err != nil {
		logfile.Fatalf("failed to create %s (%s)", JSONFile, err)
	}

	jsonWriter := bufio.NewWriter(outfile)

	defer func() {
		jsonWriter.Flush()
		outfile.Close()
	}()

	sortedLookup := make(pairSlice, 0, len(rh.machines))
	for id, index := range rh.machines {
		sortedLookup = append(sortedLookup, pair{key: index, val: id})
	}
	sort.Sort(sortedLookup)

	writeJSONFile(jsonWriter, rh, sortedLookup)
}

func writeJSONFile(w *bufio.Writer, rh *RupsHistory, sortedLookup pairSlice) {

	timestamps := fmt.Sprintf("\"timestamps\": [\n\t\"%s\"\n]", strings.Join(rh.timestamps, "\", \n\t\""))
	
	datasetObjects := make([]string, len(sortedLookup))
	for dsIndex, dsValues := range sortedLookup {
		datasetName := fmt.Sprintf("\"id\": \"%s\",", dsValues.val)
		usageArray := fmt.Sprintf("\"usages\": [\n\t%s\n]", strings.Join(rh.usages[dsValues.key], ", \n\t"))
		datasetObject := fmt.Sprintf("{%s\n%s}", datasetName, usageArray)
		datasetObjects[dsIndex] = datasetObject
	}

	json := fmt.Sprintf("{\n%s,\n\"datasets\":\n[%s]}", timestamps, strings.Join(datasetObjects, ",\n"))
	write(w, json)
}