package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/stuartleeks/wsl-clock/internal/pkg/logging"
	"github.com/stuartleeks/wsl-clock/internal/pkg/wsl"
)

// TODO
// - allow clock tolerance to be specified as an arg
// - allow max log size to be specified as an arg

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
		logging.AddMessage(err.Error())
	}
	logging.WriteLog()

	if err != nil {
		os.Exit(1)
	}
}

func execute() error {
	logging.AddMessage("********************************")
	logging.AddMessage("*** wsl-clock starting...")
	logging.AddMessage("*** Version   : %s", version)
	logging.AddMessage("*** Commit    : %s", commit)
	logging.AddMessage("*** Date      : %s", date)
	logging.AddMessage("*** Go version: %s", goversion)

	runningDistros, err := wsl.GetRunningV2Distros()
	if err != nil {
		return fmt.Errorf("Failed to get running distros: %s", err)
	}
	if len(runningDistros) == 0 {
		logging.AddMessage("No running distros - quitting")
		return nil
	}
	distroName := runningDistros[0]
	logging.AddMessage("Running commands in distro %q", distroName)

	originalTime, err := wsl.GetWslTime(distroName)
	if err != nil {
		return fmt.Errorf("Failed to get original time: %s", err)
	}

	currentTime := time.Now()
	diff := currentTime.Sub(originalTime)
	absDiffSeconds := math.Abs(diff.Seconds())

	if absDiffSeconds < toleratedTimeDiffInSeconds {
		logging.AddMessage("Time diff (%0.fs) within tolerance (%0.fs) - quitting", absDiffSeconds, toleratedTimeDiffInSeconds)
		return nil
	}

	err = wsl.ResetWslClock(distroName)
	if err != nil {
		return err
	}

	newTime, err := wsl.GetWslTime(distroName)
	if err != nil {
		return fmt.Errorf("Failed to get new time: %s", err)
	}

	logging.AddMessage("Time correction (seconds): %.0f", newTime.Sub(originalTime).Seconds())

	return nil
}
