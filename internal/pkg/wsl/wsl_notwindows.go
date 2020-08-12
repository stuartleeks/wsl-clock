// +build !windows

package wsl

import (
	"fmt"
	"time"
)

func GetRunningV2Distros() ([]string, error) {
	return []string{}, fmt.Errorf("Not implemented")
}
func GetWslTime(distroName string) (time.Time, error) {
	return time.Date(0, 0, 0, 0, 0, 0, 0, nil), fmt.Errorf("Not implemented")
}
func ResetWslClock(distroName string) error {
	return fmt.Errorf("Not implemented")
}
