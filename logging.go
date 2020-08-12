package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var message string

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
