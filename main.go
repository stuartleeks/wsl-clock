package main

import (
	"fmt"
	"math"
	"os"
	"time"
)

// Overridden via ldflags
var (
	version   = "99.0.1-devbuild"
	commit    = "unknown"
	date      = "unknown"
	goversion = "unknown"
)

const toleratedTimeDiffInSeconds float64 = 30 // allow for time executing the commands between comparisons

func main() {
	err := execute()
	if err != nil {
		addMessage(err.Error())
	}
	writeLog()

	if err != nil {
		os.Exit(1)
	}
}

func execute() error {
	addMessage("********************************")
	addMessage("*** wsl-clock starting...")
	addMessage("*** Version   : %s", version)
	addMessage("*** Commit    : %s", commit)
	addMessage("*** Date      : %s", date)
	addMessage("*** Go version: %s", goversion)

	runningDistros, err := getRunningDistros()
	if err != nil {
		return fmt.Errorf("Failed to get running distros: %s", err)
	}
	if len(runningDistros) == 0 {
		addMessage("No running distros - quitting")
		return nil
	}

	originalTime, err := getWslTime()
	if err != nil {
		return fmt.Errorf("Failed to get original time: %s", err)
	}

	currentTime := time.Now()
	diff := currentTime.Sub(originalTime)
	absDiffSeconds := math.Abs(diff.Seconds())

	if absDiffSeconds < toleratedTimeDiffInSeconds {
		addMessage("Time diff (%0.fs) within tolerance (%0.fs) - quitting", absDiffSeconds, toleratedTimeDiffInSeconds)
		return nil
	}

	err = resetWslClock()
	if err != nil {
		return err
	}

	newTime, err := getWslTime()
	if err != nil {
		return fmt.Errorf("Failed to get new time: %s", err)
	}

	addMessage("Time correction (seconds): %.0f", newTime.Sub(originalTime).Seconds())

	return nil
}
