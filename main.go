package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/alicek106/go-ec2-ssh-autoconnect/modules"
)

func printDefaultError() {
	fmt.Println(`Invalid arguments
        
        Usage: ec2-connect [command: connect, start, stop, group] [ec2-instance-name] [options]
        --key=mykey (Optional) : Use 'mykey' as a ssh private key in /etc/ec2_connect_config.ini.
                                 By default, [CONFIG][EC2_SSH_PRIVATE_KEY_DEFAULT] is used.`)
	os.Exit(100)
}

// Important!!! ###
// Handling parameter should be changed later. It is so dirty way.
func checkCommand(command string) bool {
	commands := []string{"start", "stop", "connect", "group", "list"}
	for _, value := range commands {
		if value == command {
			return true
		}
	}
	return false
}

func main() {
	// Important!!! ###
	// Handling parameter should be changed later. It is so dirty way.
	if len(os.Args) < 2 {
		printDefaultError()
	}

	if !checkCommand(os.Args[1]) {
		printDefaultError()
	}

	command := os.Args[1]
	var key string
	var instance string

	if command != "group" {
		switch {
		case len(os.Args) > 3:
			key = modules.GetEnvparser().GetCustomKey(strings.Split(os.Args[3], "=")[1])
			instance = os.Args[2]
		case len(os.Args) > 2:
			key = modules.GetEnvparser().GetDefaultKey()
			instance = os.Args[2]
		}
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
		modules.StopEc2Instances(&aem, []string{instance})
	case command == "group":
		instanceNames := modules.GetEnvparser().GetGroupInstanceNames(os.Args[3])
		if os.Args[2] == "start" {
			modules.StartEc2Instances(&aem, instanceNames)
			instanceIDs := aem.GetInstanceIDs(instanceNames)
			aem.WaitUntilActive(instanceIDs, instanceNames)
		} else if os.Args[2] == "stop" {
			modules.StopEc2Instances(&aem, instanceNames)
		}
	case command == "list":
		modules.ListEc2Instances(&aem)
	default:
		printDefaultError()
	}
}
