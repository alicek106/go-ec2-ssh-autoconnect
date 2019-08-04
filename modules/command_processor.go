package modules

import "log"

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
