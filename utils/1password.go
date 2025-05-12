package utils

import (
	"context"
	"github.com/1password/onepassword-sdk-go"
)

func ConnectOnePassword() (*onepassword.Client, error) {
	token := "TOKEN_REVOKED"

	client, err := onepassword.NewClient(
		context.TODO(),
		onepassword.WithServiceAccountToken(token),
		// TODO: Set the following to your own integration name and version.
		onepassword.WithIntegrationInfo("My 1Password Integration", "v1.0.0"),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
