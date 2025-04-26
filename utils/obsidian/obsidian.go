package obsidian

import (
	"github.com/maycon-jesus/mj-cli/utils/myIo"
	"path/filepath"
	"sync"
)

type Vault struct {
	Files []*ObsidianFile
	Notes []*ObsidianFile
	Path  string
}

func NewVault(path string) *Vault {
	return &Vault{
		Path: path,
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

			//read frontmatter
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
