package main

import (
	"fmt"
	"modules"
	"os"
)

func print_error_and_exit() {
	fmt.Println(`Invalid arguments
        
        Usage: ec2-connect [command: connect or stop] [ec2-instance-name] [options]
        --key=mykey (Optional) : Use 'mykey' as a ssh private key in /etc/ec2_connect_config.ini.
                                 By default, [CONFIG][EC2_SSH_PRIVATE_KEY_DEFAULT] is used.`)
	os.Exit(100)
}

func main() {
	if len(os.Args) < 2 {
		print_error_and_exit()
	}

	command := os.Args[1]
	var instance string
	var key string

	switch {
	case len(os.Args) > 3:
		key = os.Args[3]
		fallthrough
	case len(os.Args) > 2:
		instance = os.Args[2]
	}

    aem := modules.AwsEc2Manager{}
    aem.CheckCredentials()

	switch {
	case command == "start":
		fmt.Println("it's start")
	case command == "stop":
		fmt.Println("It's stop")
	case command == "group":
		fmt.Println("It's group")
	case command == "list":
		fmt.Println("It's list")
	default:
		print_error_and_exit()
	}

	fmt.Println(command, instance, key)
}
