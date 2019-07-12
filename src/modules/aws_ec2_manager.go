package modules

import (
	"fmt"
	"log"
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
	var err error
	if len(access_id) == 0 || len(secret_key) == 0 {
		ep := get_envparser()
		access_id, secret_key, err = ep.get_credentials()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Found credential variable in environment variables")
	}
	// TODO : Create AWS Session to check whether it is valid

	fmt.Println(access_id, secret_key)
	fmt.Println("Hello, world!")
}
