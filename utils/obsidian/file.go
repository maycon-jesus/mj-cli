package obsidian

import (
	"errors"
	"fmt"
	"github.com/maycon-jesus/mj-cli/utils/myIo"
	"maps"
	"os"
	"slices"
	"strings"
	"unicode"
)

type ObsidianFile struct {
	Name                     string            `json:"Name"`
	Path                     string            `json:"Path"`
	IsNote                   bool              `json:"IsNote"`
	Frontmatter              FilePropertiesMap `json:"Frontmatter"`
	InlineProperties         *InlineProperties `json:"InlineProperties"`
	HasBlockInlineProperties bool              `json:"HasBlockInlineProperties"`
}

func (f *ObsidianFile) GetProperty(key string) (Frontmatter *FileProperty, ok bool) {
	value, ok := f.Frontmatter[key]
	return value, ok
}

func (f *ObsidianFile) GetPropertyValues(key string) (values []string, ok bool) {
	value, ok := f.GetProperty(key)
	if ok {
		values = value.GetValues()
	}
	return values, ok
}

// AddProperty adds or updates a Frontmatter entry with the given key, Values, and metadata in the ObsidianFile's Frontmatter map.
func (f *ObsidianFile) AddProperty(key FilePropertyName, values []string) {
	value, ok := f.GetProperty(key)
	if ok {
		value.AddValues(values...)
	} else {
		frontmatter := &FileProperty{}
		frontmatter.AddValues(values...)

		f.Frontmatter[key] = frontmatter
	}
	return
}

func (f *ObsidianFile) ReadFrontmatter() *ObsidianFile {
	if !f.IsNote {
		return f
	}

	ch := myIo.ReadFile(f.Path)
	lineNumber := 0
	currentToken := ""

	//frontmatter variables temp
	lastKey := ""
	lastKeyValues := make([]string, 0)

	//inline-frontmatter variables temp
	inlineFrontmatterCount := 0
	inlineFrontmatterTag := ""
	inlineFrontmatterName := ""
	inlineFrontmatterValues := make([]string, 0)

	for text := range ch {
		lineNumber++

		// setting/removing token
		if lineNumber == 1 && text == "---" {
			currentToken = "frontmatter"
			continue
		}
		if currentToken == "frontmatter" && text == "---" {
			if lastKey != "" {
				f.AddProperty(lastKey, lastKeyValues)
			}
			currentToken = ""
			continue
		}
		if currentToken == "" && text == "%% start inline-frontmatter %%" {
			currentToken = "inline-frontmatter"
			inlineFrontmatterCount++
			f.HasBlockInlineProperties = true

			if inlineFrontmatterCount > 1 {
				panic(errors.New(fmt.Sprintf("O arquivo \"%s\" possui mais de um bloco inline-frontmatter", f.Path)))
			}
			continue
		}
		if currentToken == "inline-frontmatter" && text == "%% end inline-frontmatter %%" {
			currentToken = ""
			if inlineFrontmatterTag != "" && inlineFrontmatterName != "" {
				f.InlineProperties.AddProperty(inlineFrontmatterTag, inlineFrontmatterName, inlineFrontmatterValues)
			}
			continue
		}

		switch currentToken {
		case "frontmatter":
			text = strings.TrimSpace(text)
			isListPrefix := strings.HasPrefix(text, "-")
			if !isListPrefix {

				// Adiciona a ultima propriedade salva no map
				// Pois ja mudou o nome da propriedade
				if lastKey != "" {
					f.AddProperty(lastKey, lastKeyValues)
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
		case "inline-frontmatter":
			if text == "" {
				continue
			}

			text = strings.TrimRightFunc(text, unicode.IsSpace)
			tabsCount := strings.Count(text, "\t")
			text = strings.TrimSpace(text)
			text = strings.Replace(text, "- ", "", 1)

			switch tabsCount {
			case 0:
				inlineFrontmatterTag = text
			case 1:
				if inlineFrontmatterName != "" {
					f.InlineProperties.AddProperty(inlineFrontmatterTag, inlineFrontmatterName, inlineFrontmatterValues)
					inlineFrontmatterValues = make([]string, 0)
				}

				lineSplit := strings.SplitN(text, ":", 2)
				inlineFrontmatterName = lineSplit[0]
				if lineSplit[1] == "" {
					break
				}
				inlineFrontmatterValues = append(inlineFrontmatterValues, strings.TrimSpace(lineSplit[1]))
			case 2:
				inlineFrontmatterValues = append(inlineFrontmatterValues, text)
			}
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

	if len(f.Frontmatter) > 0 {
		fileLines = append(fileLines, "---")
	}
	for _, k := range slices.Sorted(maps.Keys(f.Frontmatter)) {
		v := f.Frontmatter[k]
		valuesLen := len(v.Values)
		if valuesLen == 0 {
			fileLines = append(fileLines, fmt.Sprintf("%s:", k))
		} else if valuesLen == 1 {
			fileLines = append(fileLines, fmt.Sprintf("%s: %s", k, v.Values[0]))
		} else {
			fileLines = append(fileLines, fmt.Sprintf("%s:", k))
			for _, value := range v.Values {
				fileLines = append(fileLines, fmt.Sprintf("  - %s", value))
			}
		}
	}
	if len(f.Frontmatter) > 0 {
		fileLines = append(fileLines, "---")
	}

	if !f.HasBlockInlineProperties {
		fileLines = append(fileLines, f.InlineProperties.ToLinesArr()...)
	}

	ignoreLines := false

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

		if text == "%% start inline-frontmatter %%" {
			ignoreLines = true
			continue
		}
		if text == "%% end inline-frontmatter %%" {
			ignoreLines = false
			fileLines = append(fileLines, f.InlineProperties.ToLinesArr()...)
			continue
		}
		if !ignoreLines {
			fileLines = append(fileLines, text)
		}
	}

	err := os.WriteFile(f.Path, []byte(strings.Join(fileLines, "\n")), 0644)
	if err != nil {
		panic(err)
	}

}

func createObsidianFile(Name string, Path string, IsNote bool) *ObsidianFile {
	return &ObsidianFile{
		Name:                     Name,
		Path:                     Path,
		IsNote:                   IsNote,
		Frontmatter:              FilePropertiesMap{},
		HasBlockInlineProperties: false,
		InlineProperties:         CreateInlineProperties(),
	}
}
