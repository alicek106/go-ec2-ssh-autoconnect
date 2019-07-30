package modules

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
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
