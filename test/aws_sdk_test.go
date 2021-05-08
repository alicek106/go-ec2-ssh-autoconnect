package test

import (
	"fmt"
	aws2 "github.com/alicek106/go-ec2-ssh-autoconnect/pkg/aws"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestDescribeEc2Instances(t *testing.T) {
	session.Must(session.NewSession())

	// Load session from shared config
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"),
	}))

	// Create new EC2 client
	ec2Svc := ec2.New(sess)
	// Call to get detailed information on each instance
	_, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		fmt.Println("TestDescribeEc2Instances: Error", err)
	} else {
		fmt.Println("TestDescribeEc2Instances: Success to describe instances")
	}
}

func TestGetUsernamePerOS(t *testing.T) {
	aem := aws2.AwsEc2Manager{Ec2StartWaitTimeout: 30}
	aem.CheckCredentials()

	username := aem.GetUsernamePerOS("bakery")
	fmt.Printf("TestGetUsernamePerOS: Found username: %s", username)
}