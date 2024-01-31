package main

import (
	"io/ioutil"
	"os"
        "testing"
        "time"

	"github.com/google/uuid"
)

const (
    DDMMYYYY = "2006-01-02"
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
	if err != nil {
		t.Fatalf("Failed to create temp uptime logs file for testing")
	}
	_, err = readUptimeLogFile(tempFilePath)
	if err == nil {
		t.Fatalf("Expected an error for corrupted uptime logs entry")
	}
        defer os.Remove(tempFilePath)
}

func TestPurgeOldUptimeLogs(t *testing.T) {
        datetime    := time.Now().UTC()
        currdate    := string((datetime.Format(DDMMYYYY)))
        olddatetime := datetime.AddDate(-1,0,0)
        olddate     := string((olddatetime.Format(DDMMYYYY)))
        PurgeOldUptimeLogs := currdate + ":000000000000001000110000\n" + olddate + ":000000000000000000010000\n"
        tempFilePath, err := createTestUptimeLogsFileWithContent(PurgeOldUptimeLogs)
        if err != nil {
                t.Fatalf("Failed to populate old uptime logs content for testing")
        }
        uptimelog, _ := readUptimeLogFile(tempFilePath)
        purgelog, _ :=  purgeOldUptimeLogs(uptimelog)
        if len(purgelog) != 1 {
                t.Fatalf("Failed to purge old uptime logs entry")
        }
        defer os.Remove(tempFilePath)
}

func TestUpdateuptimeLogs(t *testing.T) {
        datetime := time.Now().UTC()
        hour, _, _ := datetime.Clock()
        strhour := rune(hour)
        currdate := string((datetime.Format(DDMMYYYY)))
        PopulateUptimeLogs := currdate + ":000000000000000000000000\n"
        tempFilePath, err := createTestUptimeLogsFileWithContent(PopulateUptimeLogs)
        if err != nil {
                t.Fatalf("Failed to populate uptime logs content for testing")
        }
        uptimelog, _ := readUptimeLogFile(tempFilePath)
        uptimelog =  updateUptimeLogs(uptimelog)
        timeupd := string(uptimelog[currdate])
        activehr := timeupd[strhour:strhour+1]
        if activehr != "1" {
                t.Fatalf("Failed to update uptime hour ")
        }
        defer os.Remove(tempFilePath)
}
