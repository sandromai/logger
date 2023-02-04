package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type LogMessageFormatter interface {
	FormatLogMessage(
		rawLogMessage string,
	) string
}

type LogSaver interface {
	SaveLog(
		message string,
		closeProgram bool,
	)
}

type LogFileSizesChecker interface {
	CheckLogFileSizes()
}

type Logger struct {
	logsFolderPath   string
	logFileName      string
	logMessagePrefix string
}

func (logger *Logger) FormatLogMessage(
	rawLogMessage string,
) string {
	formattedLogMessage := "[" + time.Now().Format("2006-01-02 15:04:05") + "] "

	if logger.logMessagePrefix != "" {
		formattedLogMessage += logger.logMessagePrefix + " "
	}

	formattedLogMessage += rawLogMessage

	formattedLogMessage = strings.ReplaceAll(formattedLogMessage, "\r\n", "\n")
	formattedLogMessage = strings.ReplaceAll(formattedLogMessage, "\r", "\n")
	formattedLogMessage = strings.TrimSpace(formattedLogMessage)
	formattedLogMessage = strings.ReplaceAll(formattedLogMessage, "\n", " ")
	formattedLogMessage = strings.TrimSpace(formattedLogMessage)

	formattedLogMessage += "\n"

	return formattedLogMessage
}

func (logger *Logger) SaveLog(
	message string,
	closeProgram bool,
) {
	defer func() {
		if closeProgram {
			os.Exit(0)
		}
	}()

	folderStat, err := os.Stat(logger.logsFolderPath)

	if err != nil || !folderStat.IsDir() {
		return
	}

	completeLogFilePath := filepath.Join(logger.logsFolderPath, logger.logFileName)

	fileStat, err := os.Stat(completeLogFilePath)

	if err != nil || fileStat.IsDir() {
		return
	}

	fileHandler, err := os.OpenFile(completeLogFilePath, os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		return
	}

	defer fileHandler.Close()

	fileHandler.Write([]byte(
		logger.FormatLogMessage(message),
	))
}

func (logger *Logger) CheckLogFileSizes() {
	folderStat, err := os.Stat(logger.logsFolderPath)

	if err != nil || !folderStat.IsDir() {
		return
	}

	files, err := os.ReadDir(logger.logsFolderPath)

	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			return
		}

		fileName := file.Name()

		fileNameParts := strings.Split(fileName, ".")

		if fileNameParts[len(fileNameParts)-1] == "log" {
			completeLogFilePath := filepath.Join(logger.logsFolderPath, fileName)

			cmd := exec.Command("wc", "-l", completeLogFilePath)

			stdout, err := cmd.Output()

			if err != nil {
				return
			}

			logFileLineCount, err := strconv.Atoi(
				strings.SplitN(string(stdout), " ", 2)[0],
			)

			if err != nil {
				return
			}

			if logFileLineCount > 100000 {
				linesToBeRemoved := logFileLineCount - 100000

				cmd = exec.Command(
					"sed",
					"-i",
					"1,"+strconv.Itoa(linesToBeRemoved)+"d",
					completeLogFilePath,
				)

				cmd.Run()
			}
		}
	}
}

func NewLogger(logger *Logger) *Logger {
	logger.logFileName = strings.TrimSpace(logger.logFileName)

	fileNameParts := strings.Split(logger.logFileName, "/")

	logger.logFileName = fileNameParts[len(fileNameParts)-1]

	if logger.logMessagePrefix != "" {
		logger.logMessagePrefix = strings.TrimSpace(logger.logMessagePrefix)
	}

	if folderStat, err := os.Stat(logger.logsFolderPath); (err != nil && errors.Is(err, os.ErrNotExist)) || (err == nil && !folderStat.IsDir()) {
		err := os.MkdirAll(logger.logsFolderPath, 0755)

		if err != nil {
			panic("Failed to create logs folder.")
		}
	}

	completeLogFilePath := filepath.Join(logger.logsFolderPath, logger.logFileName)

	if fileStat, err := os.Stat(completeLogFilePath); (err != nil && errors.Is(err, os.ErrNotExist)) || (err == nil && fileStat.IsDir()) {
		_, err := os.Create(completeLogFilePath)

		if err != nil {
			panic("Failed to create log file.")
		}
	}

	return logger
}
