package modules

import (
	"fmt"
	"log"
	"os/exec"
	"syscall"
)

// ListEc2Instances : List all EC2 instances
func ListEc2Instances(aem *AwsEc2Manager) {
	aem.ListInstances()
}

// StartEc2Instances : Start EC2 instances
func StartEc2Instances(aem *AwsEc2Manager, instanceNames []string) {
	instanceIDs := aem.GetInstanceIDs(instanceNames)

	for index, instanceName := range instanceNames {
		log.Printf("Starting EC2 instance : %s (instance ID: %s)", instanceName, *instanceIDs[index])
	}

	aem.StartInstances(instanceIDs)
}

// ConnectSSHToInstance : Connect SSH to EC2 instance
// Don't use ssh client library or exec function! These will not work. :(
func ConnectSSHToInstance(aem *AwsEc2Manager, instanceName string, key string) {
	// TODO : Implement custom key from configuration file. (for --key=.. parameter)
	// TODO : ConnectSSHToInstance should be called in StartEc2Instances
	instanceIP := aem.GetInstancePublicIP(instanceName)
	var args = []string{"ssh", "-oStrictHostKeyChecking=no", fmt.Sprintf("ubuntu@%s", instanceIP)}
	args = append(args, fmt.Sprintf("-i%s", key))

	binary, lookErr := exec.LookPath("ssh")
	if lookErr != nil {
		panic(lookErr)
	}
	execErr := syscall.Exec(binary, args, nil)
	if execErr != nil {
		panic(execErr)
	}
}
