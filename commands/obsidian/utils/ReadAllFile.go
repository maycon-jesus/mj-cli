package utils

import (
	"bufio"
	"log"
	"os"
)

func ReadAllFile(absPath string) []string {
	file, err := os.Open(absPath)
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
	lines := make([]string, 0)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err) // Verifica erros durante a leitura
	}

	return lines
}
