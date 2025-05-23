package obsidian

import (
	"errors"
	"fmt"
	"github.com/maycon-jesus/mj-cli/utils"
	"github.com/maycon-jesus/mj-cli/utils/myIo"
	"github.com/maycon-jesus/mj-cli/utils/obsidian"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

var MonthlyCommand = &cobra.Command{
	Use:              "monthly",
	Short:            "Generate Monthly Notes",
	Aliases:          []string{"mo", "month"},
	TraverseChildren: true,
	PreRun:           CommandMonthlyValidadeFlags,
	Run:              runMonthly,
}

func CommandMonthlyValidadeFlags(cmd *cobra.Command, args []string) {
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

func GetCommandMonthly() *cobra.Command {
	todayDate := time.Now().Format("2006-01-02")
	MonthlyCommand.Flags().String("date", todayDate, "Date in format YYYY-MM-DD. If not provided, today's date will be used.\n")

	utils.CreateRequiredFlagWithViperConfig(MonthlyCommand, "template", "", "Template to use for the monthly note. If not provided, the default template will be used.", "obsidian-monthly-template-path")
	utils.CreateFlagWithViperConfig(MonthlyCommand, "dir", "", "Directory to create the monthly note. If not provided, the current directory will be used.", "obsidian-monthly-note-dir")
	MonthlyCommand.Flags().Int("quantity", 1, "Generate the next n weeks. If not provided, the current day will be used.")
	MonthlyCommand.Flags().Bool("soft", false, "Not generate error if file exists.")
	return MonthlyCommand
}

func createMonthlyFile(soft bool, date time.Time, fileContent []byte, outputDir string) {
	fileContent = []byte(obsidian.DateReplacer(string(fileContent), date))

	monthName := utils.GetMonthName(date.Month())
	monthNum := fmt.Sprintf("%02d", date.Month())

	outputPath := filepath.Join(outputDir, monthNum+" - "+monthName+".md")

	if myIo.FileExists(outputPath) {
		if soft {
			return
		}
		cobra.CheckErr(errors.New("File already exists: " + outputPath))
	}

	err := os.MkdirAll(outputDir, 0755)
	cobra.CheckErr(err)

	err = os.WriteFile(outputPath, fileContent, 0644)
	cobra.CheckErr(err)
}

func runMonthly(cmd *cobra.Command, args []string) {
	date, _ := cmd.Flags().GetString("date")
	dateTime, _ := time.Parse("2006-01-02", date)
	templatePath, _ := cmd.Flags().GetString("template")
	outputDir, _ := cmd.Flags().GetString("dir")
	nextDays, _ := cmd.Flags().GetInt("quantity")
	soft, _ := cmd.Flags().GetBool("soft")

	fileContent := []byte("")

	if templatePath != "" {
		content, err := os.ReadFile(templatePath)
		cobra.CheckErr(err)
		fileContent = content
	}

	for range nextDays {
		resolvedOutputDir := obsidian.DateReplacer(outputDir, dateTime)
		createMonthlyFile(soft, dateTime, fileContent, resolvedOutputDir)
		dateTime = dateTime.AddDate(0, 1, 0)
	}
}
