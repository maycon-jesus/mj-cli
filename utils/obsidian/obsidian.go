package obsidian

import (
	"github.com/maycon-jesus/mj-cli/utils"
	"github.com/maycon-jesus/mj-cli/utils/myIo"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Vault struct {
	Notes map[string]*ObsidianFile `json:"Notes"`
	Path  string                   `json:"Path"`
}

func NewVault(path string) *Vault {
	wd, _ := os.Getwd()
	path, _ = utils.NormalizePath(wd, path)
	return &Vault{
		Notes: make(map[string]*ObsidianFile),
		Path:  path,
	}
}

func (v *Vault) LoadAllFiles() map[string]*ObsidianFile {
	ch := make(chan myIo.File)
	wgAfterRead := &sync.WaitGroup{}
	go myIo.ListAllFilesCh(ch, v.Path)
	changedFiles := make(map[string]*ObsidianFile)

	for {
		select {
		case f, ok := <-ch:
			if !ok {
				wgAfterRead.Wait()
				return changedFiles
			}
			isNote := f.Ext == ".md"
			if !isNote {
				break
			}
			// do not if note exists without changes
			if note := v.GetNote(f.Path); note != nil && note.ModTime == f.ModTime {
				break
			}

			obsidianFile := createObsidianFile(strings.ReplaceAll(f.Name, f.Ext, ""), f.Path, isNote, f.ModTime)

			//read Frontmatter
			wgAfterRead.Add(1)
			go func() {
				defer wgAfterRead.Done()
				obsidianFile.ReadFrontmatter()
			}()

			changedFiles[obsidianFile.Path] = obsidianFile
			v.Notes[obsidianFile.Path] = obsidianFile

		}
	}
}

func (v *Vault) GetNote(path string) *ObsidianFile {
	note, _ := v.Notes[path]

	return note
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
