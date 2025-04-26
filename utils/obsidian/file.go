package obsidian

import (
	"fmt"
	"github.com/maycon-jesus/mj-cli/utils/myIo"
	"os"
	"strings"
)

type ObsidianFile struct {
	Name        string
	Path        string
	IsNote      bool
	frontmatter FilePropertiesMap
}

func (f *ObsidianFile) GetProperty(key string) (Frontmatter *FileProperty, ok bool) {
	value, ok := f.frontmatter[key]
	return value, ok
}

func (f *ObsidianFile) GetPropertyValues(key string) (values []string, ok bool) {
	value, ok := f.GetProperty(key)
	if ok {
		values = value.GetValues()
	}
	return values, ok
}

// AddProperty adds or updates a frontmatter entry with the given key, values, and metadata in the ObsidianFile's frontmatter map.
func (f *ObsidianFile) AddProperty(key FilePropertyName, values []string, metadata FilePropertyMetadata) {
	value, ok := f.GetProperty(key)
	if ok {
		value.AddValues(values...)
		value.AddMetadataMap(metadata)
	} else {
		frontmatter := &FileProperty{}
		frontmatter.AddValues(values...)
		frontmatter.SetMetadata(metadata)

		f.frontmatter[key] = frontmatter
	}
	return
}

func (f *ObsidianFile) ReadFrontmatter() *ObsidianFile {
	if !f.IsNote {
		return f
	}

	ch := myIo.ReadFile(f.Path)
	ignoreLines := false
	lineNumber := 0
	lastKey := ""
	lastKeyValues := make([]string, 0)

	for text := range ch {
		if ignoreLines {
			continue
		}

		lineNumber++
		if lineNumber == 1 && text != "---" {
			ignoreLines = true
			continue
		}
		if lineNumber == 1 && text == "---" {
			continue
		}

		if text == "---" {
			if lastKey != "" {
				f.AddProperty(lastKey, lastKeyValues, FilePropertyMetadata{})
			}
			ignoreLines = true
			continue
		}

		isListPrefix := strings.HasPrefix(text, "-")
		if !isListPrefix {

			// Adiciona a ultima propriedade salva no map
			// Pois ja mudou o nome da propriedade
			if lastKey != "" {
				f.AddProperty(lastKey, lastKeyValues, FilePropertyMetadata{})
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

	for k, v := range f.frontmatter {
		if !strings.HasPrefix(k, "metadata.") {
			continue
		}

		metadataParts := strings.Split(k, ".")
		if len(metadataParts) != 3 {
			continue
		}

		targetPropertyName := metadataParts[1]

		metadataPropertyName := metadataParts[2]

		targetProperty, ok := f.GetProperty(targetPropertyName)
		if !ok {
			continue
		}
		if len(v.values) == 0 {
			targetProperty.AddMetadata(metadataPropertyName, "")
		} else {
			targetProperty.AddMetadata(metadataPropertyName, v.values[0])
		}
	}

	return f
}

func (f *ObsidianFile) WriteFile() {
	if !f.IsNote {
		return
	}

	ch := myIo.ReadFile(f.Path)
	frontmatterLimiterCount := 0
	hasFrontmatter := false
	lineCounter := 0
	fileLines := make([]string, 0)

	if len(f.frontmatter) > 0 {
		fileLines = append(fileLines, "---")
	}
	for k, v := range f.frontmatter {
		valuesLen := len(v.values)
		if valuesLen == 0 {
			fileLines = append(fileLines, fmt.Sprintf("%s:", k))
		} else if valuesLen == 1 {
			fileLines = append(fileLines, fmt.Sprintf("%s: %s", k, v.values[0]))
		} else {
			fileLines = append(fileLines, fmt.Sprintf("%s:", k))
			for _, value := range v.values {
				fileLines = append(fileLines, fmt.Sprintf("  - %s", value))
			}
		}
	}
	if len(f.frontmatter) > 0 {
		fileLines = append(fileLines, "---")
	}

	for text := range ch {
		lineCounter++
		if lineCounter == 1 && text == "---" {
			frontmatterLimiterCount++
			hasFrontmatter = true
			continue
		}
		if hasFrontmatter && frontmatterLimiterCount < 2 {
			if text == "---" {
				frontmatterLimiterCount++
			}
			continue
		}
		fileLines = append(fileLines, text)
	}

	err := os.WriteFile(f.Path, []byte(strings.Join(fileLines, "\n")), 0644)
	if err != nil {
		panic(err)
	}

}

func createObsidianFile(Name string, Path string, IsNote bool) *ObsidianFile {
	return &ObsidianFile{
		Name:        Name,
		Path:        Path,
		IsNote:      IsNote,
		frontmatter: FilePropertiesMap{},
	}
}
