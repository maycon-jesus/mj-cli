package myIo

import (
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	cobra.CheckErr(err)
	return info.IsDir()
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") {
			return false
		}
		cobra.CheckErr(err)
	}
	return !info.IsDir()
}
