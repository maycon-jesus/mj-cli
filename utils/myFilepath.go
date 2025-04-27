package utils

import (
	"errors"
	"path/filepath"
)

func NormalizePath(globalPath string, path string) (string, error) {
	if !filepath.IsAbs(globalPath) {
		return "", errors.New("globalPath must be abslute path")
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(globalPath, path)
	}
	path = filepath.Clean(path)

	return path, nil
}
