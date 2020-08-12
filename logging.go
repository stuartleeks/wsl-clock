package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var message string

const maxLogFileSize int64 = 5 * 1024 * 1024

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
	backupLogPath := filepath.Join(userProfile, ".wsl-clock.old.log")

	handleLogFileRotation(logPath, backupLogPath)

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
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

func handleLogFileRotation(logPath string, backupLogPath string) {
	size, err := getFileSize(logPath)
	if err != nil {
		if err != nil {
			fmt.Printf("Error getting log file size %q: %s", logPath, err)
			panic(err)
		}
	}
	if size > maxLogFileSize {
		if _, err = os.Stat(backupLogPath); err != nil {
			if !os.IsNotExist(err) {
				fmt.Printf("Error checking backup log path %q: %s", backupLogPath, err)
				panic(err)
			}
		} else {
			os.Remove(backupLogPath)
		}
		os.Rename(logPath, backupLogPath)
	}
}

func getFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil // treat not found as empty for this use-case
		}
		return 0, err
	}
	return info.Size(), nil
}
