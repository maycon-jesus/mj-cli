package logger

import (
	"os"
	"sync"
)

type FileWriter struct {
	file *os.File
	mu   sync.Mutex
}

func NewFileWriter(file *os.File) *FileWriter {
	return &FileWriter{file: file, mu: sync.Mutex{}}
}

func (fw *FileWriter) Write(p []byte) (n int, err error) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	return fw.file.Write(p)
}

func (fw *FileWriter) Close() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	return fw.file.Close()
}
