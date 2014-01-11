package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	HtmlFile = "html/index.html"
)

func (rh *RupsHistory) createHtmlPage() {
	outfile, err := os.Create(HtmlFile)
	if err != nil {
		logfile.Fatalf("failed to create %s (%s)", HtmlFile, err)
	}

	htmlWriter := bufio.NewWriter(outfile)

	defer func() {
		htmlWriter.Flush()
		outfile.Close()
	}()

	sortedLookup := make(pairSlice, 0, len(rh.machines))
	for id, index := range rh.machines {
		sortedLookup = append(sortedLookup, pair{key: index, val: id})
	}
	sort.Sort(sortedLookup)

	write(htmlWriter, Header)

	writeDatasets(htmlWriter, rh, sortedLookup)

	writePlots(htmlWriter, rh, sortedLookup)

	write(htmlWriter, Middle)

	writePlotDivs(htmlWriter, rh, sortedLookup)

	write(htmlWriter, Footer)
}

func writePlotDivs(w *bufio.Writer, rh *RupsHistory, sortedLookup pairSlice) {
	for dsIndex, _ := range sortedLookup {
		div := fmt.Sprintf(`<div class="demo-container"> <div id="placeholder%d" class="demo-placeholder"></div> </div>%s`, dsIndex, "\n")
		write(w, div)
	}
}

func writePlots(w *bufio.Writer, rh *RupsHistory, sortedLookup pairSlice) {
	for dsIndex, dsValues := range sortedLookup {
		plotStart := fmt.Sprintf(`$.plot("#placeholder%d", [`, dsIndex)
		write(w, plotStart)
		dsName := fmt.Sprintf("d%d", dsIndex)
		str := fmt.Sprintf("{ label: \"%s\", data: %s }\n", dsValues.val, dsName)
		write(w, str)
		write(w, PlotClose)
	}
}

func writeDatasets(w *bufio.Writer, rh *RupsHistory, sortedLookup pairSlice) {
	for dsIndex, dsValues := range sortedLookup {
		dsName := fmt.Sprintf("d%d", dsIndex)

		values := fmt.Sprintf("%s = [ ", dsName)
		for dpIndex, dpValue := range rh.usages[dsValues.key] {
			values = fmt.Sprintf("%s[moment(\"%s\").toDate(), %s], \n", values, rh.timestamps[dpIndex], dpValue)
		}
		values = fmt.Sprintf("%s];\n", values)

		write(w, values)
	}
}
