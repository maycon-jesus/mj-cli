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

var MonthCommand = &cobra.Command{
	Use:              "month",
	Short:            "Utilitários para o obsidian",
	TraverseChildren: true,
	Run:              runMonth,
}

func GetCommandMonth() *cobra.Command {
	return MonthCommand
}

func runMonth(cmd *cobra.Command, args []string) {
	vaultDir, _ := cmd.Flags().GetString("vault-dir")
	journalDir := viper.GetString("obsidianBulletJournalDir")

	notePath := filepath.Join(vaultDir, journalDir)

	timeNow := time.Now()
	year, _ := timeNow.ISOWeek()
	month := int(timeNow.Month())
	monthName := ""

	switch month {
	case 1:
		monthName = "Janeiro"
	case 2:
		monthName = "Fevereiro"
	case 3:
		monthName = "Março"
	case 4:
		monthName = "Abril"
	case 5:
		monthName = "Maio"
	case 6:
		monthName = "Junho"
	case 7:
		monthName = "Julho"
	case 8:
		monthName = "Agosto"
	case 9:
		monthName = "Setembro"
	case 10:
		monthName = "Outubro"
	case 11:
		monthName = "Novembro"
	case 12:
		monthName = "Dezembro"
	}

	//adiciona ano no path
	notePath = filepath.Join(notePath, strconv.Itoa(year))

	//adiciona caminho do mes no path
	weeklyDir := viper.GetString("obsidianBulletJournalMonthDir")
	notePath = filepath.Join(notePath, weeklyDir)

	//adiciona arquivo do mes no path
	filename := fmt.Sprintf("%02d - %s.md", month, monthName)
	notePath = filepath.Join(notePath, filename)

	_, err := os.Stat(notePath)
	if err == nil {
		fmt.Println("Ja existe um arquivo para este mês!")
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
