package logger

import (
	"log"
	"os"
)

// Setup initializes the logger (placeholder for more advanced logging)
func Setup() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Info(v ...interface{}) {
	log.Println("[INFO]", v)
}

func Error(v ...interface{}) {
	log.Println("[ERROR]", v)
}