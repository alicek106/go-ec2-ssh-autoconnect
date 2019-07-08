package modules

import (
	"fmt"
	"os"
)

type AwsEc2Manager struct {
	session string
	client  string
}

func (aem *AwsEc2Manager) CheckCredentials() {
	// TODO : Check Credential from configuration or env var
	access_id := os.Getenv("AWS_ACCESS_KEY_ID")
	secret_key := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if len(access_id) == 0 || len(secret_key) == 0 {
		fmt.Println("Cannot find credential variable in environment variables")
		// TODO : Find credential in configuration (using EnvParser)
	} else {
		fmt.Println("Found credential variable in environment variables")
	}

	fmt.Println(access_id, secret_key)
	fmt.Println("Hello, world!")
}
