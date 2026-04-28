package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode/utf8"
)

type LogEntry struct {
	File    string
	Line    int
	IsValid bool
	Reason  string
}

func rSplit(s, sep string) (before, after string) {
	if sep == "" {
		return s, ""
	}
	i := strings.LastIndex(s, sep)
	if i >= 0 {
		return s[:i], s[i+len(sep):]
	}
	return s, ""
}

func validateAndStructure(filetype string, line string) LogEntry {
	res := LogEntry{IsValid: false}
	if !utf8.ValidString(line) {
		res.Reason = "Invalid UTF-8 encoding"
		return res
	}

	switch filetype {
	case "auth.log", "secure":
		if strings.Contains(line, "sshd") || strings.Contains(line, "sudo") {
			res.IsValid = true
			res.Reason = "Matched sshd/sudo"
		} else {
			res.Reason = "Missing sshd/sudo"
		}
	case "audit.log":
		if strings.Contains(line, "type=SYSCALL") {
			res.IsValid = true
			res.Reason = "Matched type=SYSCALL"
		} else {
			res.Reason = "Missing type=SYSCALL"
		}
	case "syslog", "kern.log", "messages":
		res.IsValid = true
	}
	return res
}

func getFileType(path string) string {
	_, after := rSplit(path, "/")
	return after
}

func readLogs(path string, out chan<- LogEntry, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.Open(path)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to open log file: %s\n", err)
		if err != nil {
			return
		}
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	filetype := getFileType(path)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		out <- validateAndStructure(filetype, scanner.Text())
	}

}

func main() {
	// Test script to get used to Go :)
	// Debian, RHEL/CentOS/Fedora log journals
	// NOTE: other distros (e.g. Arch Linux) use systemd's `journald` for auth/system logs
	logsToParse := []string{
		"/var/log/auth.log",
		"/var/log/secure",
		"/var/log/syslog",
		"/var/log/messages",
		"/var/log/kern.log",
		"/var/log/audit/audit.log",
	}
	resChan := make(chan LogEntry)
	var wg sync.WaitGroup

	for i := range logsToParse {
		wg.Add(1)
		go readLogs(logsToParse[i], resChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resChan)
	}()

	for res := range resChan {
		if res.IsValid {
			fmt.Printf("[OK] %s: %s\n", res.File, res.Reason)
		} else {
			fmt.Printf("[WARN] %s: %s\n", res.File, res.Reason)
		}
	}
}
