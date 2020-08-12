package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var message string

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
	// TODO - add version
	addMessage("********************************")
	addMessage("*** Update WSL clock starting...")

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

func addMessage(newMessage string, a ...interface{}) {
	if message != "" {
		message += "\n"
	}
	timestamp := time.Now().UTC().Format("2006-01-02 15:04:05 ")
	message += timestamp + fmt.Sprintf(newMessage, a...)
}
func writeLog() {
	userProfile := os.Getenv("USERPROFILE")
	logPath := filepath.Join(userProfile, ".wsl-clock.log")

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		fmt.Printf("Error opening log file %q: %s", logPath, err)
		panic(err)
	}
	defer file.Close()
	_, err = file.WriteString(message + "\n")
	if err != nil {
		fmt.Printf("Error writing to log file %q: %s", logPath, err)
		panic(err)
	}
}

func getRunningDistros() ([]string, error) {
	// TODO - consider changing this to us verbose listing and test for any running v2 instances
	// and then use that to determine an instance to run the remaining WSL commands in
	output, err := execCmdToLines("wsl.exe", "--list", "--running", "--quiet")
	if err != nil {
		return []string{}, err
	}
	return output, nil
}
func getWslTime() (time.Time, error) {
	output, err := execCmdToLines("wsl.exe", "sh", "-c", "date -Iseconds")
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed to call WSL to get current time: %s", err)
	}

	timeString := output[0]
	timeValue, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed to parse time %q: %s", timeString, err)
	}

	return timeValue, nil
}
func resetWslClock() error {
	_, err := execCmdToLines("wsl.exe", "-u", "root", "sh", "-c", "hwclock -s")
	if err != nil {
		return fmt.Errorf("Failed to call WSL to reset clock: %s", err)
	}
	return nil
}

func execCmdToLines(program string, arg ...string) ([]string, error) {
	cmd := exec.Command(program, arg...)
	outputTemp, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	output := outputTemp
	if len(outputTemp) >= 2 && outputTemp[1] == 0 {
		output = make([]byte, len(outputTemp)/2)
		for i := 0; i < len(output); i++ {
			output[i] = outputTemp[2*i]
		}
	}

	reader := bytes.NewReader(output)
	scanner := bufio.NewScanner(reader)
	if scanner == nil {
		return []string{}, fmt.Errorf("Failed to parse stdout")
	}
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return []string{}, fmt.Errorf("Error reading stdout: %s", err)
	}

	return lines, nil
}
