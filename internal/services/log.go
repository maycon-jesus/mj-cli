package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
)

type LogService struct {
	logger *slog.Logger
}

func NewLogService() *LogService {
	return &LogService{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func (s *LogService) WithLogger(logger *slog.Logger) *LogService {
	s.logger = logger
	return s
}

type LogLine struct {
	Error error
	Log   map[string]any
}

func (s *LogService) GetLogs(filePath string) (chan LogLine, error) {
	s.logger.Debug("Getting logs from file", "filePath", filePath)

	ch := make(chan LogLine)

	if !path.IsAbs(filePath) {
		wd, err := os.Getwd()
		if err != nil {
			s.logger.Debug("Failed to get working directory", "error", err)
			return ch, err
		}
		filePath = path.Join(wd, filePath)
		s.logger.Debug("Resolved log file path to absolute path", "filePath", filePath)
	}

	s.logger.Debug("Checking if log file exists", "filePath", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		s.logger.Debug("Failed to open log file", "filePath", filePath, "error", err)
		return ch, err
	}

	go func() {
		s.logger.Debug("Started log file reading goroutine", "filePath", filePath)
		defer s.logger.Debug("Finished log file reading goroutine", "filePath", filePath)
		defer file.Close()
		defer close(ch)
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Bytes()
			var lineData map[string]any
			err := json.Unmarshal(line, &lineData)
			if err != nil {
				s.logger.Debug("Failed to unmarshal log line", "error", err)
				ch <- LogLine{Error: fmt.Errorf("error parsing log line: %w", err), Log: nil}
				continue
			}
			ch <- LogLine{Error: nil, Log: lineData}
		}
		if err := scanner.Err(); err != nil {
			s.logger.Debug("Error reading log file", "filePath", filePath, "error", err)
			ch <- LogLine{Error: fmt.Errorf("error reading log file: %w", err)}
		}
	}()

	return ch, nil
}
