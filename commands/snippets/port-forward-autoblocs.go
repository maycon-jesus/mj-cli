package snippets

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"time"
)

// PortForwardAutoblocsApiSnippetCommand defines a Cobra command to manage port-forwarding for the Autoblocs API service.
var PortForwardAutoblocsApiSnippetCommand = &cobra.Command{
	Use: "port-forward-autoblocs-api",
	Run: RunPortForwardAutoblocsApiSnippetCommand,
}

// GetPortForwardAutoblocsApiSnippetCommand returns the Cobra command for managing port forwarding in Autoblocs API.
func GetPortForwardAutoblocsApiSnippetCommand() *cobra.Command {
	return PortForwardAutoblocsApiSnippetCommand
}

// runCommandRealtime executes a shell command in real-time, streaming stdout and stderr based on the visibility setting.
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

// runCommand executes a shell command with the given arguments and returns its combined output and error, if any.
func runCommand(command string, args []string) (string, error) {
	cmde := exec.Command(command, args...)
	content, err := cmde.CombinedOutput()
	return string(content), err
}

// checkIsLogged validates if the user has an active session by attempting to retrieve an identity token using gcloud CLI.
// Returns true if successfully authenticated, otherwise returns false.
func checkIsLogged() bool {
	_, err := runCommand("gcloud", []string{"auth", "print-identity-token"})

	if err != nil {
		return false
	}
	return true
}

// gcloudLogin ensures the user is authenticated with gcloud by checking logged-in status and triggering login if necessary.
func gcloudLogin() {
	logged := checkIsLogged()
	if logged {
		return
	}

	err := runCommandRealtime("gcloud", []string{"auth", "login"}, false)
	cobra.CheckErr(err)
}

// RunPortForwardAutoblocsApiSnippetCommand establishes a port-forwarding session for the "autoblocs-api" service in a loop.
// It ensures the user is authenticated with gcloud and uses kubectl to forward local port 8081 to the remote service port 8080.
func RunPortForwardAutoblocsApiSnippetCommand(cmd *cobra.Command, args []string) {
	//client, err := utils.ConnectOnePassword()
	//cobra.CheckErr(err)

	gcloudLogin()

	for {
		runCommandRealtime("kubectl", []string{"port-forward", "--namespace", "apps", "autoblocs-api-infra-api-0", "8081:8080"}, false)
		time.Sleep(3 * time.Second)
	}
}
