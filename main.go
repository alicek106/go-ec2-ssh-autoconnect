package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alicek106/go-ec2-ssh-autoconnect/modules"
)

func printDefaultError() {
	fmt.Println(`Invalid arguments
        
        Usage: ec2-connect [command: connect or stop] [ec2-instance-name] [options]
        --key=mykey (Optional) : Use 'mykey' as a ssh private key in /etc/ec2_connect_config.ini.
                                 By default, [CONFIG][EC2_SSH_PRIVATE_KEY_DEFAULT] is used.`)
	os.Exit(100)
}

func main() {
	if len(os.Args) < 2 {
		printDefaultError()
	}

	command := os.Args[1]
	instance := os.Args[2]
	var key string

	switch {
	case len(os.Args) > 3:
		key = modules.GetEnvparser().GetCustomKey(strings.Split(os.Args[3], "=")[1])
	case len(os.Args) > 2:
		key = modules.GetEnvparser().GetDefaultKey()
	}

	// TODO: Ec2StartWaitTimeout should be able to set by CLI parameter, later :D
	aem := modules.AwsEc2Manager{Ec2StartWaitTimeout: 30}
	aem.CheckCredentials()

	switch {
	case command == "connect":
		modules.ConnectSSHToInstance(&aem, instance, key)
	case command == "start":
		modules.StartEc2Instances(&aem, []string{instance})
	case command == "stop":
		fmt.Println("It's stop")
	case command == "group":
		fmt.Println("It's group")
	case command == "list":
		modules.ListEc2Instances(&aem)
	default:
		printDefaultError()
	}
}
