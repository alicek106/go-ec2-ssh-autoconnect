package modules

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Ec2InstanceInfo : Data struct for storing EC2 instance data
type Ec2InstanceInfo struct {
	instancdID   string
	instanceName string
	status       string
	publicIP     string
}

// AwsEc2Manager : AWS EC2 Session and client amanger
type AwsEc2Manager struct {
	session             *session.Session
	client              *ec2.EC2
	Ec2StartWaitTimeout int
}

// CheckCredentials : Check existing credential from shell or configuartion
func (aem *AwsEc2Manager) CheckCredentials() {
	// TODO : Check Credential from configuration or env var
	accessID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var err error
	if len(accessID) == 0 || len(secretKey) == 0 {
		log.Printf("Cannot found credential in environment variable.")
		ep := GetEnvparser()
		accessID, secretKey, err = ep.GetCredentials()
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Found credential in configuration file.")
		}
	} else {
		log.Println("Found credential variable in environment variables")
	}

	aem.ValidateCredential(accessID, secretKey)
}

// ValidateCredential : Validate AWS Credential
func (aem *AwsEc2Manager) ValidateCredential(accessID string, secretKey string) {
	session.Must(session.NewSession())

	// Load session
	aem.session = session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-2"),
		Credentials: credentials.NewStaticCredentials(accessID, secretKey, ""),
	}))

	// Test AWS function using provided credential
	aem.client = ec2.New(aem.session)
	_, err := aem.client.DescribeInstances(nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Succeed to validate AWS credential.")
	}
}

// GetInstanceIDs : Return instance IDs from instanceNames
func (aem *AwsEc2Manager) GetInstanceIDs(instanceNames []string) []*string {
	var instanceIDs = []*string{}
	for _, instanceName := range instanceNames {
		// If filters is defined directly in input parameter, it triggers syntax error.
		filter := aem.getFilterForName(instanceName)
		result, err := aem.client.DescribeInstances(filter)
		if err != nil {
			fmt.Println(err)
		} else {
			instanceIDs = append(instanceIDs, result.Reservations[0].Instances[0].InstanceId)
		}
	}
	return instanceIDs
}

// ListInstances : List multiple instances.
func (aem *AwsEc2Manager) ListInstances() {
	result, err := aem.client.DescribeInstances(nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// var ec2NameMaxLength int
		var ec2InstanceList = []Ec2InstanceInfo{}
		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				var ec2InstanceInfo = Ec2InstanceInfo{
					status:     *instance.State.Name,
					instancdID: *instance.InstanceId,
				}

				if instance.PublicIpAddress == nil {
					ec2InstanceInfo.publicIP = "Unknown"
				} else {
					ec2InstanceInfo.publicIP = *instance.PublicIpAddress
				}

				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
						ec2InstanceInfo.instanceName = *tag.Value
					}
				}
				ec2InstanceList = append(ec2InstanceList, ec2InstanceInfo)
			}
		}

		writer := tabwriter.NewWriter(os.Stdout, 16, 8, 2, '\t', 0)
		fmt.Fprintln(writer, "Instance ID\tInstance Name\tIP Address\tStatus")
		for _, ec2InstanceInfo := range ec2InstanceList {
			formatting := fmt.Sprintf("%s\t%s\t%s\t%s", ec2InstanceInfo.instancdID,
				ec2InstanceInfo.instanceName, ec2InstanceInfo.publicIP, ec2InstanceInfo.status)
			fmt.Fprintln(writer, formatting)
		}

		writer.Flush()
	}
}

// StartInstances : Start multiple instances.
func (aem *AwsEc2Manager) StartInstances(instanceIDs []*string) {
	input := &ec2.StartInstancesInput{
		InstanceIds: instanceIDs, // It should be used with pointer
		DryRun:      aws.Bool(true),
	}
	_, err := aem.client.StartInstances(input)
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		_, err := aem.client.StartInstances(input)
		if err != nil {
			log.Fatal("Error", err)
		} else {
			log.Printf("Succeed to start EC2 instances.")
		}
	} else { // This could be due to a lack of permissions
		log.Fatal("Error", err)
	}
}

// WaitUntilActive : Wait unil all instances are up and running.
// Order of instanceIDs and instanceNames is same. (Refer to GetInstanceIDs)
func (aem *AwsEc2Manager) WaitUntilActive(instanceIDs []*string, instanceNames []string) {
	for index := range instanceIDs {
		log.Printf("Start to waiting EC2 instance %s (instance ID : %s)...", instanceNames[index], *instanceIDs[index])
		for tries := 1; tries <= aem.Ec2StartWaitTimeout; tries++ {
			// If EC2 instance is not running state
			if aem.getInstanceStatus(instanceIDs[index:index+1]) != "running" {
				// If EC2 instance is not running state after 30s, continue to next instance.
				if tries == 30 {
					log.Printf("Failed to wait for EC2 instance to be active.")
					break
				} else {
					// Wait seconds (aem.Ec2StartWaitTimeout)
					log.Printf("Waiting for starting EC2 instance.. %d tries.", tries)
					time.Sleep(time.Second)
				}
			} else {
				// EC2 instance is already running, pass to wait 30 seconds for warm-up
				if tries == 1 {
					log.Printf("EC2 instance is in active.")
					break
				}
				// Wait 30 seconds until SSH daemon is ready
				log.Printf("EC2 instance is in active. Waiting for 30 seconds for warm-up.")
				time.Sleep(30 * time.Second)
				break
			}
		}
	}
}

// GetInstancePublicIP : Return EC2 instance public IP from instance name
func (aem *AwsEc2Manager) GetInstancePublicIP(instanceName string) (publicIP string) {
	filter := aem.getFilterForName(instanceName)
	result, err := aem.client.DescribeInstances(filter)
	if err != nil {
		fmt.Println(err)
	}
	return *result.Reservations[0].Instances[0].PublicIpAddress
}

// StopInstances : Stop multiple instances.
func (aem *AwsEc2Manager) StopInstances(instanceIDs []*string) {
	input := &ec2.StopInstancesInput{
		InstanceIds: instanceIDs, // It should be used with pointer
		DryRun:      aws.Bool(true),
	}
	_, err := aem.client.StopInstances(input)
	awsErr, ok := err.(awserr.Error)

	if ok && awsErr.Code() == "DryRunOperation" {
		input.DryRun = aws.Bool(false)
		_, err := aem.client.StopInstances(input)
		if err != nil {
			log.Fatal("Error", err)
		} else {
			log.Printf("Succeed to stop EC2 instances.")
		}
	} else { // This could be due to a lack of permissions
		log.Fatal("Error", err)
	}
}

func (aem *AwsEc2Manager) getInstanceStatus(instanceID []*string) (status string) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: instanceID,
	}
	result, err := aem.client.DescribeInstances(input)
	if err != nil {
		log.Println("Error in waiting for EC2 instances to be active")
		log.Fatal(err)
	} else {
		status = *result.Reservations[0].Instances[0].State.Name
	}
	return status
}

// getFilterForName : Return filter for describing instances
func (aem *AwsEc2Manager) getFilterForName(instanceName string) (input *ec2.DescribeInstancesInput) {
	filters := []*ec2.Filter{
		&ec2.Filter{
			Name:   aws.String("tag:Name"),
			Values: []*string{aws.String(instanceName)},
		},
	}
	input = &ec2.DescribeInstancesInput{
		Filters: filters,
	}
	return input
}
