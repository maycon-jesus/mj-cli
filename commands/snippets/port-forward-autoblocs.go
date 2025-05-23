package snippets

import (
	"github.com/maycon-jesus/mj-cli/utils/terminal"
	"github.com/spf13/cobra"
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

// checkIsLogged validates if the user has an active session by attempting to retrieve an identity token using gcloud CLI.
// Returns true if successfully authenticated, otherwise returns false.
func checkIsLogged() bool {
	_, err := terminal.RunCommand("gcloud", []string{"auth", "print-identity-token"})

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

	err := terminal.RunCommandRealtime("gcloud auth login", terminal.RunCommandOptions{})
	cobra.CheckErr(err)
}

// RunPortForwardAutoblocsApiSnippetCommand establishes a port-forwarding session for the "autoblocs-api" service in a loop.
// It ensures the user is authenticated with gcloud and uses kubectl to forward local port 8081 to the remote service port 8080.
func RunPortForwardAutoblocsApiSnippetCommand(cmd *cobra.Command, args []string) {
	//client, err := utils.ConnectOnePassword()
	//cobra.CheckErr(err)

	gcloudLogin()

	for {
		terminal.RunCommandRealtime("kubectl port-forward --namespace apps autoblocs-api-infra-api-0 8081:8080", terminal.RunCommandOptions{})
		time.Sleep(3 * time.Second)
	}
}
