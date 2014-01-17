package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var (
	logfile          = log.New(os.Stderr, "", log.LstdFlags)
	ScpDestination   = flag.String("scpremote", "", "The destination for the scp command. You need to set up keys for this.")
	MaxHistoryPoints = flag.Int("numpts", 50, "The maximum number of history points to maintain.")
	Frequency        = flag.Int("freq", 5, "The frequency at which to poll rups. {1, 5, 15}")
)

func main() {
	flag.Parse()

	if *Frequency != 1 && *Frequency != 5 && *Frequency != 15 {
		panic("Frequency provide is not 1, 5 or 15")
	}

	history := readRupsHistory()

	quit := make(chan bool)
	done := make(chan bool)

	//Handle cleanup and saving when this process gets killed
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		select {
		case <-c:
			quit <- true
		}
	}()

	//Update at frequency minute intervals
	ticker := time.NewTicker(time.Duration(*Frequency) * time.Minute)
	update(history)
	go func() {
		for {
			select {
			case <-ticker.C:
				update(history)
			case <-quit:
				ticker.Stop()
				done <- true
				return
			}
		}
	}()

	<-done

}

func update(rh *RupsHistory) {
	rh.addCurrentRups()
	rh.createJSONFile()

	if *ScpDestination != "" {
		scpCommand := fmt.Sprintf("scp %s %s", JSONFile, *ScpDestination)
		c := exec.Command("/bin/sh", "-c", scpCommand)
		output, err := c.CombinedOutput()
		if err != nil {
			logfile.Fatalf("failed output=[%s]: %s\n", output, err)
		}
	}
	rh.saveRupsHistory()
}
