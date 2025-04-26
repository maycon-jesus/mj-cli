package obsidian

import (
	"github.com/maycon-jesus/mj-cli/utils"
	"github.com/maycon-jesus/mj-cli/utils/myIo"
	"os"
	"path/filepath"
	"sync"
)

type Vault struct {
	Files       []*ObsidianFile
	Notes       []*ObsidianFile
	Path        string
	TagsDirPath string
	Tag
}

func NewVault(path string, tagsDirPath string) *Vault {
	wd, _ := os.Getwd()
	path, _ = utils.NormalizePath(wd, path)
	tagsDirPath, _ = utils.NormalizePath(path, tagsDirPath)
	return &Vault{
		Path:        path,
		TagsDirPath: tagsDirPath,
	}
}

func (v *Vault) LoadAllFiles() {
	ch := make(chan myIo.File)
	wgAfterRead := &sync.WaitGroup{}
	go myIo.ListAllFilesCh(ch, v.Path)

	for {
		select {
		case f, ok := <-ch:
			if !ok {
				wgAfterRead.Wait()
				return
			}
			isNote := f.Ext == ".md"
			obsidianFile := createObsidianFile(f.Name, f.Path, isNote)

			//read Frontmatter
			wgAfterRead.Add(1)
			go func() {
				defer wgAfterRead.Done()
				obsidianFile.ReadFrontmatter()
			}()

			v.Files = append(v.Files, obsidianFile)
			if isNote {
				v.Notes = append(v.Notes, obsidianFile)
			}
		}
	}
}

func (v *Vault) GetNote(path string) *ObsidianFile {
	var npath string
	if filepath.IsAbs(path) {
		npath = path
	} else {
		npath = filepath.Join(v.Path, path)
	}
	for _, note := range v.Notes {
		if note.Path == npath {
			return note
		}
	}

	return nil
}

func (v *Vault) GetDirectoryNotes(path string) []*ObsidianFile {
	var npath string
	if filepath.IsAbs(path) {
		npath = path
	} else {
		npath = filepath.Join(v.Path, path)
	}
	var notes []*ObsidianFile
	for _, note := range v.Notes {
		if filepath.Dir(note.Path) == npath {
			notes = append(notes, note)
		}
	}

	return notes
}

func (v *Vault) GetTagTemplateNote(tag string) (*ObsidianFile, bool) {
	path := filepath.Join(v.TagsDirPath, tag+".md")

	for _, note := range v.Notes {
		if note.Path == path {
			return note, true
		}
	}

	return nil, false
}
