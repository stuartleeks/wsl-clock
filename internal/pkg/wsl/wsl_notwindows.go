// +build !windows

package wsl

import (
	"fmt"
	"time"
)

func GetRunningDistros() ([]string, error) {
	return []string{}, fmt.Errorf("Not implemented")
}
func GetWslTime() (time.Time, error) {
	return time.Date(0, 0, 0, 0, 0, 0, 0, nil), fmt.Errorf("Not implemented")
}
func ResetWslClock() error {
	return fmt.Errorf("Not implemented")
}
