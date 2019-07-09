package test

import (
	"fmt"
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
	result, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("Success", result)
	}
}
