package modules

import (
	"fmt"
	"log"
	"os"
)

// AwsEc2Manager : AWS EC2 Session and client amanger
type AwsEc2Manager struct {
	session string
	client  string
}

// CheckCredentials : Check existing credential from shell or configuartion
func (aem *AwsEc2Manager) CheckCredentials() {
	// TODO : Check Credential from configuration or env var
	accessID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var err error
	if len(accessID) == 0 || len(secretKey) == 0 {
		ep := getEnvparser()
		accessID, secretKey, err = ep.getCredentials()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Found credential variable in environment variables")
	}
	// TODO : Create AWS Session to check whether it is valid

	fmt.Println(accessID, secretKey)
	fmt.Println("Hello, world!")
}
