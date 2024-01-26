package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/uuid"
)

func createTestUptimeLogsFileWithContent(content string) (string, error) {
	tempFile, err := ioutil.TempFile("", "testUptimeLogs")
	if err != nil {
		return "", err
	}
	tempFilePath := tempFile.Name()
	if _, err := tempFile.WriteString(content); err != nil {
		os.Remove(tempFilePath)
		return "", err
	}
	if err := tempFile.Close(); err != nil {
		os.Remove(tempFilePath)
		return "", err
	}
	return tempFilePath, nil
}

func TestUptimeLogFileDoesNotExist(t *testing.T) {
	bogusUptimeLogsFilePath := uuid.New().String()
	uptimeLogs, err := readUptimeLogFile(bogusUptimeLogsFilePath)
	if uptimeLogs != nil || err != nil {
		t.Fatalf("Expected err and uptimeLogs to be nil if uptime logs file does not exist")
	}
}

func TestCorruptedUptimeLogs(t *testing.T) {
	corruptedUptimeLogs := `2024-01-18:000000000000001000110000
2024-01-13000000000000000000010000`
	tempFilePath, err := createTestUptimeLogsFileWithContent(corruptedUptimeLogs)
	defer os.Remove(tempFilePath)
	if err != nil {
		t.Fatalf("Failed to create temp uptime logs file for testing")
	}
	_, err = readUptimeLogFile(tempFilePath)
	if err == nil {
		t.Fatalf("Expected an error for currupted uptime logs entry")
	}
}
