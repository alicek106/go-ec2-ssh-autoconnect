package modules

// ListEc2Instances : List all EC2 instances
func ListEc2Instances(aem *AwsEc2Manager) {
	aem.ListInstances()
}

// StartEc2Instances : Start EC2 instances
func StartEc2Instances(aem *AwsEc2Manager, instanceNames []string) {
	instanceIDs := aem.GetInstanceIDs(instanceNames)
	aem.StartInstances(instanceIDs)
	// aem.WaitUntilActive(instanceNames)
}
