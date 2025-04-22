package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Frontmatter = map[string][]string

type FrontmatterMetadata struct {
	order int
}

func ReadFrontmatter(c chan<- Frontmatter, absPath string) {
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

	lineNumber := 0
	lastKey := ""
	lastKeyValues := make([]string, 0)
	allFrontmatter := make(Frontmatter)

	for scanner.Scan() {
		lineNumber++
		text := strings.TrimSpace(scanner.Text())
		if lineNumber == 1 && text != "---" {
			break
		}
		if lineNumber == 1 && text == "---" {
			continue
		}
		if text == "---" {
			if lastKey != "" {
				allFrontmatter[lastKey] = lastKeyValues
			}
			break
		}

		isListPrefix := strings.HasPrefix(text, "-")
		if !isListPrefix {

			// Adiciona a ultima propriedade salva no map
			// Pois ja mudou o nome da propriedade
			if lastKey != "" {
				allFrontmatter[lastKey] = lastKeyValues
				lastKeyValues = nil
				lastKey = ""
			}

			frontmatter := strings.SplitN(text, ":", 2)
			frontmatterKey := strings.TrimSpace(frontmatter[0])
			frontmatterValue := strings.TrimSpace(frontmatter[1])
			lastKey = frontmatterKey

			if frontmatterValue != "" {
				lastKeyValues = append(lastKeyValues, frontmatterValue)
			}
		} else {
			lastKeyValues = append(lastKeyValues, strings.ReplaceAll(text, "- ", ""))
		}

	}

	c <- allFrontmatter

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
