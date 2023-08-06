package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func handler() {
	destroy()
}

func main() {
	lambda.Start(handler)
}

// EC2DescribeInstancesAPI defines the interface for the DescribeInstances function.
// We use this interface to test the function using a mocked service.
type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}
type EC2TerminanteInstanceAPI interface {
	TerminateInstances(ctx context.Context,
		params *ec2.TerminateInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}

func getInstances(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}

func terminateInstances(ctx context.Context, api EC2TerminanteInstanceAPI, parameter4Destroy *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	return api.TerminateInstances(ctx, parameter4Destroy)
}
func destroy() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{}

	result, err := getInstances(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error retrieving information about your Amazon EC2 instances:")
		fmt.Println(err)
		return
	}

	terminate := &ec2.TerminateInstancesInput{
		InstanceIds: []string{},
	}

	for _, r := range result.Reservations {
		fmt.Println("Reservation ID: " + *r.ReservationId)
		for _, i := range r.Instances {
			fmt.Println("   " + *i.InstanceId)

			if *i.InstanceId != "i-06c23e4f25705c1c9" {
				fmt.Sprintf("Instance IDs: %s", *i.InstanceId)
				terminate.InstanceIds = append(terminate.InstanceIds, *i.InstanceId) //add all instances to this slice
			}
		}
		//call api to terminate all instance just 1 time with all instances
		terminateInstances(context.TODO(), client, terminate)

		fmt.Println("")
	}
}
