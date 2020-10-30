package logging

import (
	"fmt"
	"log"
)

func LogInfo(v ...interface{}) {
	printLn(fmt.Sprintf("INFO: %v", v))
}
func LogDebug(v ...interface{}) {
	printLn(fmt.Sprintf("DEBUG: %v", v))
}
func LogError(v ...interface{}) {
	printLn(fmt.Sprintf("ERROR: %v", v))
}

func printLn(message string) {
	log.Println(message)
}
