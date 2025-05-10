package obsidian

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var WeekCommand = &cobra.Command{
	Use:              "week",
	Short:            "Utilitários para o obsidian",
	Aliases:          []string{"wk"},
	TraverseChildren: true,
	Run:              runWeek,
}

func GetCommandWeek() *cobra.Command {
	return WeekCommand
}

func runWeek(cmd *cobra.Command, args []string) {
	vaultDir, _ := cmd.Flags().GetString("vault-dir")
	journalDir := viper.GetString("obsidianBulletJournalDir")

	notePath := filepath.Join(vaultDir, journalDir)

	timeNow := time.Now()
	year, weekNum := timeNow.ISOWeek()

	//adiciona ano no path
	notePath = filepath.Join(notePath, strconv.Itoa(year))

	//adiciona caminho da semana no path
	weeklyDir := viper.GetString("obsidianBulletJournalWeeklyDir")
	notePath = filepath.Join(notePath, weeklyDir)

	//adiciona arquivo da semana no path
	filename := fmt.Sprintf("Semana %02d.md", weekNum)
	notePath = filepath.Join(notePath, filename)

	_, err := os.Stat(notePath)
	if err == nil {
		fmt.Println("Ja existe um arquivo para esta semana!")
		return
	}

	if !errors.Is(err, os.ErrNotExist) {
		fmt.Println(err)
		return
	}

	fmt.Println("Arquivo criado com sucesso")
	os.MkdirAll(filepath.Dir(notePath), 0755)
	os.WriteFile(notePath, []byte(""), 0644)
}
