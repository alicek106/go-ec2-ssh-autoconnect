package modules

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

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
	session *session.Session
	client  *ec2.EC2
}

// CheckCredentials : Check existing credential from shell or configuartion
func (aem *AwsEc2Manager) CheckCredentials() {
	// TODO : Check Credential from configuration or env var
	accessID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	var err error
	if len(accessID) == 0 || len(secretKey) == 0 {
		log.Printf("Cannot found credential in environment variable.")
		ep := getEnvparser()
		accessID, secretKey, err = ep.getCredentials()
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
		log.Printf("Success to validate AWS credential.")
	}
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
						// if len(*tag.Value) > ec2NameMaxLength {
						// 	ec2NameMaxLength = len(*tag.Value)
						// }
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
	log.Printf("Starting EC2 instance...")
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
			log.Printf("Succeed to start all instances. Waiting for instances to be active..")
		}
	} else { // This could be due to a lack of permissions
		log.Fatal("Error", err)
	}
}

// WaitUntilActive : Wait unil all instances are up and running.
func (aem *AwsEc2Manager) WaitUntilActive(instanceIDs []*string) {
	// TODO : Implement for waiting EC2 instances to be running
}
