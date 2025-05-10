package obsidian

import (
	"errors"
	"github.com/maycon-jesus/mj-cli/utils"
	"github.com/maycon-jesus/mj-cli/utils/myIo"
	"github.com/maycon-jesus/mj-cli/utils/obsidian"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

var DailyCommand = &cobra.Command{
	Use:              "daily",
	Short:            "Generate Daily Notes",
	Args:             cobra.NoArgs,
	PreRun:           CommandDailyValidadeFlags,
	Run:              runDaily,
	TraverseChildren: true,
}

func CommandDailyValidadeFlags(cmd *cobra.Command, args []string) {
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
	outputDirResolved := obsidian.DateReplacer(outputDirNormalized, dateTime)
	err = cmd.Flags().Lookup("dir").Value.Set(outputDirResolved)
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

func GetCommandDaily() *cobra.Command {
	todayDate := time.Now().Format("2006-01-02")
	DailyCommand.Flags().String("date", todayDate, "Date in format YYYY-MM-DD. If not provided, today's date will be used.\n")

	utils.CreateRequiredFlagWithViperConfig(DailyCommand, "template", "", "Template to use for the daily note. If not provided, the default template will be used.", "obsidian-daily-template-path")
	utils.CreateFlagWithViperConfig(DailyCommand, "dir", "", "Directory to create the daily note. If not provided, the current directory will be used.", "obsidian-daily-note-dir")
	DailyCommand.Flags().Int("quantity", 1, "Generate the next n days. If not provided, the current day will be used.")

	return DailyCommand
}

func createDailyFile(date time.Time, fileContent []byte, outputDir string) {
	fileContent = []byte(obsidian.DateReplacer(string(fileContent), date))

	outputPath := filepath.Join(outputDir, date.Format("2006-01-02")+".md")
	if myIo.FileExists(outputPath) {
		cobra.CheckErr(errors.New("File already exists: " + outputPath))
	}

	err := os.MkdirAll(outputDir, 0755)
	cobra.CheckErr(err)

	err = os.WriteFile(outputPath, fileContent, 0644)
	cobra.CheckErr(err)
}

func runDaily(cmd *cobra.Command, args []string) {
	date, _ := cmd.Flags().GetString("date")
	dateTime, _ := time.Parse("2006-01-02", date)
	templatePath, _ := cmd.Flags().GetString("template")
	outputDir, _ := cmd.Flags().GetString("dir")
	nextDays, _ := cmd.Flags().GetInt("quantity")

	fileContent := []byte("")

	if templatePath != "" {
		content, err := os.ReadFile(templatePath)
		cobra.CheckErr(err)
		fileContent = content
	}

	for range nextDays {
		createDailyFile(dateTime, fileContent, outputDir)
		dateTime = dateTime.AddDate(0, 0, 1)
	}

}
