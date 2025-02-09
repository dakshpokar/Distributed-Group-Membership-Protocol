package logger

import (
	"fmt"
	"os"
)

var File *os.File

func InitializeLog() {
	fileName := "vm.log"

	// Open the file, create it if it doesn't exist, and truncate it if it does
	var err error
	File, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}

}

func LogToFile(message string) {
	// Write the message to the file
	_, err := File.WriteString(message)
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		return
	}
}
