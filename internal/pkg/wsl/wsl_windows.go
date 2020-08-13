// +build windows

package wsl

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

func GetRunningV2Distros() ([]string, error) {
	// TODO - consider changing this to us verbose listing and test for any running v2 instances
	// and then use that to determine an instance to run the remaining WSL commands in
	output, err := execCmd("wsl.exe", "--list", "--verbose")
	if err != nil {
		return []string{}, err
	}

	distros, err := parseDistroOutput(string(output))
	if err != nil {
		return []string{}, err
	}

	v2Distros := []string{}
	for _, distro := range distros {
		if distro.Version == "2" && distro.State == "Running" {
			v2Distros = append(v2Distros, distro.Name)
		}
	}
	return v2Distros, nil
}
func GetWslTime(distroName string) (time.Time, error) {
	output, err := execCmdToLines("wsl.exe", "-d", distroName, "sh", "-c", "date -Iseconds")
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
func ResetWslClock(distroName string) error {
	_, err := execCmdToLines("wsl.exe", "-d", distroName, "-u", "root", "sh", "-c", "hwclock -s")
	if err != nil {
		return fmt.Errorf("Failed to call WSL to reset clock: %s", err)
	}
	return nil
}

func execCmd(program string, arg ...string) ([]byte, error) {
	cmd := exec.Command(program, arg...)

	const CREATE_NO_WINDOW = 0x08000000
	sysAttr := syscall.SysProcAttr{}
	sysAttr.CreationFlags = CREATE_NO_WINDOW
	sysAttr.HideWindow = true
	cmd.SysProcAttr = &sysAttr

	outputTemp, err := cmd.Output()
	if err != nil {
		return []byte{}, err
	}

	output := outputTemp
	if len(outputTemp) >= 2 && outputTemp[1] == 0 {
		output = make([]byte, len(outputTemp)/2)
		for i := 0; i < len(output); i++ {
			output[i] = outputTemp[2*i]
		}
	}

	return output, nil
}

func execCmdToLines(program string, arg ...string) ([]string, error) {

	output, err := execCmd(program, arg...)
	if err != nil {
		return []string{}, err
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
