package main

import (
	"os"
	"path/filepath"
)

func main() {
	current_path, err := os.Executable()

	if err != nil {
		return
	}

	logger := NewLogger(&Logger{
		logsFolderPath:   filepath.Join(filepath.Dir(current_path), "logs"),
		logFileName:      "test.log",
		logMessagePrefix: "[TEST]:",
	})

	logger.CheckLogFileSizes()

	logger.SaveLog("Testing!", false)
}
