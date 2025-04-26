package myIo

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type File struct {
	Name string
	Path string
	Ext  string
}

func ListAllFiles(absPath string) []*File {
	var files []*File
	ch := make(chan File)
	go ListAllFilesCh(ch, absPath)

	for {
		select {
		case f, ok := <-ch:
			if !ok {
				return files
			}
			files = append(files, &f)
		}
	}
}

func ListAllFilesCh(ch chan<- File, absPath string) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		wg.Wait()
		close(ch)
	}()

	go readDirectoryFiles(ch, wg, absPath)
}

func readDirectoryFiles(ch chan<- File, wg *sync.WaitGroup, absPath string) {
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
			go readDirectoryFiles(ch, wg, filepath.Join(absPath, entry.Name()))
		} else {
			//if filepath.Ext(entry.Name()) != ".md" {
			//	continue
			//}
			f := File{
				Path: filepath.Join(absPath, entry.Name()),
				Name: entry.Name(),
				Ext:  filepath.Ext(entry.Name()),
			}

			ch <- f
		}
	}
}
