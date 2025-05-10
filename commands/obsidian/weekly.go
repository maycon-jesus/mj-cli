package obsidian

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/maycon-jesus/mj-cli/utils"
	"github.com/maycon-jesus/mj-cli/utils/myIo"
	"github.com/maycon-jesus/mj-cli/utils/obsidian"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var WeeklyCommand = &cobra.Command{
	Use:              "weekly",
	Short:            "Generate Weekly Notes",
	Aliases:          []string{"wk"},
	TraverseChildren: true,
	PreRun:           CommandWeeklyValidadeFlags,
	Run:              runWeekly,
}

func CommandWeeklyValidadeFlags(cmd *cobra.Command, args []string) {
	vaultDir, err := cmd.Flags().GetString("vault-dir")
	cobra.CheckErr(err)

	//date
	date, err := cmd.Flags().GetString("date")
	cobra.CheckErr(err)
	dateTime, err := time.Parse("2006-01-02", date)
	cobra.CheckErr(err)

	//dir
	outputDir, err := cmd.Flags().GetString("dir")
	cobra.CheckErr(err)
	outputDirNormalized, err := utils.NormalizePath(vaultDir, outputDir)
	cobra.CheckErr(err)
	err = cmd.Flags().Lookup("dir").Value.Set(outputDirNormalized)
	cobra.CheckErr(err)

	//template
	templatePath, err := cmd.Flags().GetString("template")
	cobra.CheckErr(err)
	if templatePath != "" {
		templatePathNormalized, err := utils.NormalizePath(vaultDir, templatePath)
		cobra.CheckErr(err)
		templatePathResolved := obsidian.DateReplacer(templatePathNormalized, dateTime)
		templatePathExists := myIo.FileExists(templatePathResolved)
		if !templatePathExists {
			cobra.CheckErr(errors.New("Template does not exist: " + templatePathResolved))
		}
		err = cmd.Flags().Lookup("template").Value.Set(templatePathResolved)
	}
}

func GetCommandWeekly() *cobra.Command {
	todayDate := time.Now().Format("2006-01-02")
	WeeklyCommand.Flags().String("date", todayDate, "Date in format YYYY-MM-DD. If not provided, today's date will be used.\n")

	utils.CreateRequiredFlagWithViperConfig(WeeklyCommand, "template", "", "Template to use for the weekly note. If not provided, the default template will be used.", "obsidian-weekly-template-path")
	utils.CreateFlagWithViperConfig(WeeklyCommand, "dir", "", "Directory to create the weekly note. If not provided, the current directory will be used.", "obsidian-weekly-note-dir")
	WeeklyCommand.Flags().Int("quantity", 1, "Generate the next n weeks. If not provided, the current day will be used.")
	return WeeklyCommand
}

func createWeeklyFile(date time.Time, fileContent []byte, outputDir string) {
	fileContent = []byte(obsidian.DateReplacer(string(fileContent), date))

	_, weekN := date.ISOWeek()

	outputPath := filepath.Join(outputDir, "Semana "+strconv.Itoa(weekN)+".md")
	fmt.Println(outputPath)
	if myIo.FileExists(outputPath) {
		cobra.CheckErr(errors.New("File already exists: " + outputPath))
	}

	err := os.MkdirAll(outputDir, 0755)
	cobra.CheckErr(err)

	err = os.WriteFile(outputPath, fileContent, 0644)
	cobra.CheckErr(err)
}

func runWeekly(cmd *cobra.Command, args []string) {
	date, _ := cmd.Flags().GetString("date")
	dateTime, _ := time.Parse("2006-01-02", date)
	templatePath, _ := cmd.Flags().GetString("template")
	outputDir, _ := cmd.Flags().GetString("dir")
	nextDays, _ := cmd.Flags().GetInt("quantity")

	fmt.Println(outputDir)

	fileContent := []byte("")

	if templatePath != "" {
		content, err := os.ReadFile(templatePath)
		cobra.CheckErr(err)
		fileContent = content
	}

	correctionFactor := int(dateTime.Weekday()) * -1
	dateTime = dateTime.AddDate(0, 0, correctionFactor)

	for range nextDays {
		resolvedOutputDir := obsidian.DateReplacer(outputDir, dateTime)
		weeklyDates := ""
		for i := 0; i < 7; i++ {
			if i > 0 {
				weeklyDates += " | "
			}
			weeklyDates += fmt.Sprintf("[[../03 - Diários/%s|%s]]", dateTime.Format("2006-01-02"), dateTime.Format("02/01"))
			dateTime = dateTime.AddDate(0, 0, 1)
		}

		nFileContent := bytes.ReplaceAll(fileContent, []byte("{{WEEKLY_DATES}}"), []byte(weeklyDates))
		createWeeklyFile(dateTime, nFileContent, resolvedOutputDir)
	}
}
