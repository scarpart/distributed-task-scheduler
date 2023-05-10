package logger

import (
	"log"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrLogger     *log.Logger
)
//
//func init() {
//	// The file path below should be relative, but for some reason it's returning errors even though the path should be correct.
//	// I'm leaving the entire path here for now, just to make it work, then I'm going to change things around a bit. 
//	file, err := os.OpenFile("/home/artzmeister/code/portfolio/distributed-task-scheduler/logs/api_logs.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
//	if err != nil {
//		log.Fatal("Could not set up logger:", err)
//		return
//	}
//	WarningLogger = log.New(file, "[WARN] - ", log.Ldate|log.Ltime|log.Lshortfile)
//	InfoLogger = log.New(file, "[INFO] - ", log.Ldate|log.Ltime|log.Lshortfile)
//	ErrLogger = log.New(file, "[ERR] - ", log.Ldate|log.Ltime|log.Lshortfile)
//}
