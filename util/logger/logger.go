package logger

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrLogger     *log.Logger
)

func init() {
	file, err := os.OpenFile("logs/api_logs.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Could not set up logger:", err)
		return
	}
	WarningLogger = log.New(file, "[WARN] - ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(file, "[INFO] - ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrLogger = log.New(file, "[ERR] - ", log.Ldate|log.Ltime|log.Lshortfile)
}
