package myIo

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func readFile(ch chan<- string, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		text := strings.TrimSpace(scanner.Text())
		ch <- text
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	close(ch)
}

func ReadFile(path string) <-chan string {
	ch := make(chan string)
	go readFile(ch, path)
	return ch
}
