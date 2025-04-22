package obsidian

import (
	"fmt"
	"github.com/maycon-jesus/mj-cli/commands/obsidian/utils"
	"github.com/maycon-jesus/mj-cli/utils/mySlices"
	"github.com/spf13/cobra"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

var TagsProperties = &cobra.Command{
	Use:   "tags-properties",
	Short: "tags-properties short",
	Run:   run,
}

type file struct {
	absPath     string
	name        string
	frontmatter utils.Frontmatter
	lines       []string
}

func (f *file) setFrontmatter(frontmatter utils.Frontmatter) {
	f.frontmatter = frontmatter
}

func (f *file) addFrontmatterField(fieldName string, fieldValues ...string) {
	values := make([]string, 0)
	for _, value := range fieldValues {
		values = append(values, value)
	}
	f.frontmatter[fieldName] = values
}

func (f *file) readFile() {
	f.lines = utils.ReadAllFile(f.absPath)
}

func (f *file) updateLinesFrontmatter() {
	fileLines := make([]string, len(f.lines))
	copy(fileLines, f.lines)

	for i, line := range fileLines {
		if i == 0 && line != "---" {
			break
		}
		if i != 0 && line == "---" {
			fileLines = fileLines[i+1:]
		}
	}

	propertiesLines := make([]string, 0)

	if len(f.frontmatter) != 0 {
		propertiesLines = append(propertiesLines, "---")
	}

	for _, property := range slices.Sorted(maps.Keys(f.frontmatter)) {
		values := f.frontmatter[property]
		if len(values) == 0 {
			propertiesLines = append(propertiesLines, fmt.Sprintf("%s:", property))
		} else if len(values) == 1 {
			propertiesLines = append(propertiesLines, fmt.Sprintf("%s: %s", property, values[0]))
		} else {
			propertiesLines = append(propertiesLines, fmt.Sprintf("%s:", property))
			for _, value := range values {
				propertiesLines = append(propertiesLines, fmt.Sprintf("  - %s", value))
			}
		}
	}

	//for property, values := range f.frontmatter {
	//	if len(values) <= 1 {
	//		propertiesLines = append(propertiesLines, fmt.Sprintf("%s: %s", property, values[0]))
	//	} else {
	//		propertiesLines = append(propertiesLines, fmt.Sprintf("%s:", property))
	//		for _, value := range values {
	//			propertiesLines = append(propertiesLines, fmt.Sprintf("  - %s", value))
	//		}
	//	}
	//}

	if len(f.frontmatter) != 0 {
		propertiesLines = append(propertiesLines, "---")
	}

	fileLinesWithProperties := make([]string, 0, len(fileLines)+len(propertiesLines))

	for _, line := range propertiesLines {
		fileLinesWithProperties = append(fileLinesWithProperties, line)
	}
	for _, line := range fileLines {
		fileLinesWithProperties = append(fileLinesWithProperties, line)
	}
	err := os.WriteFile(f.absPath, []byte(strings.Join(fileLinesWithProperties, "\n")), 0644)
	if err != nil {
		panic(err)
	}

}

func GetCommandTagsProperties() *cobra.Command {
	TagsProperties.Flags().StringP("templates-dir", "t", "99 - Meta/02 - Tags", "tags templates directory")
	return TagsProperties
}

func run(cmd *cobra.Command, args []string) {
	vaultDir, _ := cmd.Flags().GetString("vault-dir")
	templatesDir, _ := cmd.Flags().GetString("templates-dir")
	var vaultDirAbs string
	var templatesDirAbs string

	if !filepath.IsAbs(vaultDir) {
		wdDir, _ := os.Getwd()
		vaultDirAbs = filepath.Join(wdDir, vaultDir)
	} else {
		vaultDirAbs = vaultDir
	}

	if !filepath.IsAbs(templatesDir) {
		templatesDirAbs = filepath.Join(vaultDirAbs, templatesDir)
	} else {
		templatesDirAbs = templatesDir
	}
	fmt.Println(templatesDirAbs)

	files := readAllFiles(vaultDirAbs)

	for _, file := range files {
		if filepath.Dir(file.absPath) == templatesDirAbs {
			continue
		}
		fmt.Println("===============")
		fmt.Println("Arquivo atual:", file.absPath)
		fmt.Println(file.frontmatter)
		fileUpdated := false
		tags, ok := file.frontmatter["tags"]
		if !ok {
			continue
		}
		for _, tag := range tags {
			fmt.Println("Tag atual:", tag)
			tagTemplatePath := filepath.Join(templatesDirAbs, tag+".md")
			for _, fileTemplate := range files {
				if fileTemplate.absPath != tagTemplatePath {
					continue
				}

				frontmatterKeys := mySlices.MapKeysToSlice[string, []string](fileTemplate.frontmatter)
				frontmatterKeysFiltered := mySlices.Filter[string](frontmatterKeys, func(key string) bool {
					if strings.HasPrefix(key, "metadata.") {
						return false
					}
					return true
				})

				for _, propertieTemplate := range frontmatterKeysFiltered {
					if _, ok := file.frontmatter[propertieTemplate]; ok == false {
						fmt.Println("O arquivo não possui a propriedade:", propertieTemplate)
						file.addFrontmatterField(propertieTemplate, "")
						fileUpdated = true
					}
				}

			}
		}

		if fileUpdated {
			file.readFile()
			file.updateLinesFrontmatter()
		}
	}

}

func readAllFiles(absPath string) []*file {
	var files []*file
	ch := make(chan file)
	wg := &sync.WaitGroup{}

	go func() {
		wg.Wait()
		close(ch)
	}()

	wg.Add(1)
	go listAllFiles(ch, wg, absPath)

	for {
		select {
		case f, ok := <-ch:
			if !ok {
				return files
			}
			go readFileFrontmatter(wg, &f)
			files = append(files, &f)
		}
	}
}

func listAllFiles(ch chan<- file, wg *sync.WaitGroup, absPath string) {
	defer func() {
		wg.Done()
	}()
	dirEntry, err := os.ReadDir(absPath)
	if err != nil {
		fmt.Printf("Erro ao ler diretório: %v\n", err)
		return
	}

	for _, entry := range dirEntry {
		if entry.IsDir() {
			wg.Add(1)
			go listAllFiles(ch, wg, filepath.Join(absPath, entry.Name()))
		} else {
			if filepath.Ext(entry.Name()) != ".md" {
				continue
			}
			f := file{
				absPath:     filepath.Join(absPath, entry.Name()),
				name:        entry.Name(),
				frontmatter: nil,
			}

			//Adicionado 1 pois para cada arquivo precisa ler também o frontmatter
			wg.Add(1)
			ch <- f
		}
	}
}

func readFileFrontmatter(wg *sync.WaitGroup, f *file) {
	defer wg.Done()
	ch := make(chan utils.Frontmatter)
	go utils.ReadFrontmatter(ch, f.absPath)
	frontmatter := <-ch
	f.setFrontmatter(frontmatter)
}
