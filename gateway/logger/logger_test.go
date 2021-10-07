package logger_test

import (
	"encoding/csv"
	"net/http"
	"os"
	"testing"

	"github.com/future-architect/apidoor/gateway/logger"
)

func TestUpdateLog(t *testing.T) {
	// open file
	file, err := os.OpenFile(os.Getenv("LOG_PATH"), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		t.Fatal(err)
	}

	// save current environment variable
	tmp := os.Getenv("LOG_PATTERN")
	t.Cleanup(func() {
		os.Setenv("LOG_PATTERN", tmp)
	})

	// run test
	r, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	r.Header.Set("TEST1", "header1")
	r.Header.Set("TEST2", "header2")
	logger.LogOptionPattern = []logger.LogOption{
		logger.WithTime(),
		logger.WithKey(),
		logger.WithPath(),
		logger.HeaderElement("TEST1"),
		logger.HeaderElement("TEST2"),
	}
	for i := 0; i < 2; i++ {
		logger.UpdateLog("key", "path", r)
	}

	// check if log is valid
	reader := csv.NewReader(file)
	recordNum := 0

	for {
		line, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Fatal(err)
		}
		if line[1] != "key" {
			t.Fatalf("unexpected log %s, expected 'key'", line[1])
		} else if line[2] != "path" {
			t.Fatalf("unexpected log %s, expected 'path'", line[2])
		} else if line[3] != "header1" {
			t.Fatalf("unexpected log %s, expected 'header1'", line[3])
		} else if line[4] != "header2" {
			t.Fatalf("unexpected log %s, expected 'header2'", line[4])
		}
		recordNum++
	}

	if recordNum != 2 {
		t.Fatalf("unexpected number of log %d, expected 2", recordNum)
	}
}