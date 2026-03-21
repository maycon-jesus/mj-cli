package filesystem

import (
	"os"
	"path"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetFileContent(path string) (string, error) {
	exists := Exists(path)
	if !exists {
		return "", os.ErrNotExist
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func WriteFile(path string, content string) error {
	stat, err := os.Stat(path)
	if err == nil && stat.IsDir() {
		return os.ErrInvalid
	}
	return os.WriteFile(path, []byte(content), stat.Mode().Perm())
}

func GetAbsolutePath(p string) (string, error) {
	isAbs := path.IsAbs(p)
	if !isAbs {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		p = path.Join(wd, p)
	}
	return p, nil
}
