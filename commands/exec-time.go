package commands

import (
	"fmt"
	"github.com/maycon-jesus/mj-cli/utils"
	"github.com/maycon-jesus/mj-cli/utils/myIo"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
	"time"
)

var ExecTimeCommand = &cobra.Command{
	Use:    "exec-time",
	Short:  "Return execution time of a snippet",
	Run:    RunExecTimeCommand,
	PreRun: ValidadeFlagsExecTimeCommand,
}

func GetExecTimeCommand() *cobra.Command {
	ExecTimeCommand.Flags().String("wd", "", "Working directory")
	return ExecTimeCommand
}

func runCommandRealtime(workDirectory string, command string, args []string) error {
	cmde := exec.Command(command, args...)
	cmde.Stdout = os.Stdout
	cmde.Stderr = os.Stderr
	cmde.Dir = workDirectory
	err := cmde.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}

func ValidadeFlagsExecTimeCommand(cmd *cobra.Command, args []string) {
	flagWd, err := cmd.Flags().GetString("wd")
	cobra.CheckErr(err)
	wd, _ := os.Getwd()
	if flagWd == "" {
		cmd.Flags().Lookup("wd").Value.Set(wd)
	} else {
		pathNormalized, _ := utils.NormalizePath(wd, flagWd)
		cmd.Flags().Lookup("wd").Value.Set(pathNormalized)
	}
	flagWd, err = cmd.Flags().GetString("wd")
	cobra.CheckErr(err)
	wdExists := myIo.DirectoryExists(flagWd)
	if !wdExists {
		cobra.CheckErr(fmt.Errorf("%s is not a directory", flagWd))
	}
}

// RunExecTimeCommand executes a shell command, measures its execution time, and displays the time in various units.
func RunExecTimeCommand(cmd *cobra.Command, args []string) {
	wd, _ := cmd.Flags().GetString("wd")
	cmdParsed := strings.Split(args[0], " ")
	bin := cmdParsed[0]
	argsParsed := cmdParsed[1:]
	startedAt := time.Now()
	err := runCommandRealtime(wd, bin, argsParsed)
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(startedAt)
	fmt.Println("\n\n====[EXECUTION TIME]====")
	fmt.Printf("Milliseconds: %d\n", elapsed.Milliseconds())
	fmt.Printf("Seconds: %.0f\n", elapsed.Seconds())
	fmt.Printf("Minutes: %.0f\n", elapsed.Minutes())
	fmt.Println("========================")
}
