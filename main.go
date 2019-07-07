package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`Invalid arguments
        
        Usage: ec2-connect [command: connect or stop] [ec2-instance-name] [options]
        --key=mykey (Optional) : Use 'mykey' as a ssh private key in /etc/ec2_connect_config.ini.
                                 By default, [CONFIG][EC2_SSH_PRIVATE_KEY_DEFAULT] is used.`)
		return
	}

	command := os.Args[1]
	instance := os.Args[2]
	key := os.Args[3]
	fmt.Println(command, instance, key)
	fmt.Println("Hello, world!")
}
