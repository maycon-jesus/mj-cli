package snippets

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"time"
)

var PortForwardAutoblocsApiSnippetCommand = &cobra.Command{
	Use: "port-forward-autoblocs-api",
	Run: RunPortForwardAutoblocsApiSnippetCommand,
}

func GetPortForwardAutoblocsApiSnippetCommand() *cobra.Command {
	return PortForwardAutoblocsApiSnippetCommand
}

func runCommandRealtime(command string, args []string, hideContent bool) error {
	cmde := exec.Command(command, args...)
	printer := MyPrinter{hideContent: hideContent}
	cmde.Stdout = printer
	cmde.Stderr = printer
	err := cmde.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}

func runCommand(command string, args []string) (string, error) {
	cmde := exec.Command(command, args...)
	content, err := cmde.CombinedOutput()
	return string(content), err
}

func checkIsLogged() bool {
	_, err := runCommand("gcloud", []string{"auth", "print-identity-token"})

	if err != nil {
		return false
	}
	return true
}

func gcloudLogin() {
	logged := checkIsLogged()
	if logged {
		return
	}

	err := runCommandRealtime("gcloud", []string{"auth", "login"}, false)
	cobra.CheckErr(err)
}

func RunPortForwardAutoblocsApiSnippetCommand(cmd *cobra.Command, args []string) {
	//client, err := utils.ConnectOnePassword()
	//cobra.CheckErr(err)

	gcloudLogin()

	for successExec := false; successExec == false; {
		err := runCommandRealtime("kubectl", []string{"port-forward", "--namespace", "apps", "autoblocs-api-infra-api-0", "8081:8080"}, false)
		if err == nil {
			successExec = true
		} else {
			time.Sleep(3 * time.Second)
		}
	}
}
