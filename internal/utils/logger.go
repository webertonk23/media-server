package utils
import (
	"bufio"
	"io"
	"log"
	"os"
	"sync"
	"time"
)
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}
type MemoryLogger struct {
	entries []LogEntry
	maxSize int
	mu      sync.RWMutex
	multi   io.Writer
}
var (
	GlobalLogger *MemoryLogger
	logFilePath  string
)
func InitLogger(maxSize int, filePath string) {
	logFilePath = filePath
	GlobalLogger = &MemoryLogger{
		entries: make([]LogEntry, 0, maxSize),
		maxSize: maxSize,
	}
	var writers []io.Writer
	writers = append(writers, os.Stdout, GlobalLogger)
	if filePath != "" {
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			writers = append(writers, f)
		}
	}
	GlobalLogger.multi = io.MultiWriter(writers...)
	log.SetOutput(GlobalLogger.multi)
}
func (l *MemoryLogger) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	msg := string(p)
	entry := LogEntry{
		Timestamp: time.Now(),
		Message:   msg,
	}
	if len(l.entries) >= l.maxSize {
		l.entries = l.entries[1:]
	}
	l.entries = append(l.entries, entry)
	return len(p), nil
}
func (l *MemoryLogger) GetLogs() []LogEntry {
	if logFilePath != "" {
		return readLogsFromFile(logFilePath, l.maxSize)
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	res := make([]LogEntry, len(l.entries))
	copy(res, l.entries)
	return res
}
func readLogsFromFile(path string, n int) []LogEntry {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > n {
			lines = lines[1:]
		}
	}
	entries := make([]LogEntry, len(lines))
	for i, line := range lines {
		entries[i] = LogEntry{
			Timestamp: time.Now(), 
			Message:   line,
		}
	}
	return entries
}
