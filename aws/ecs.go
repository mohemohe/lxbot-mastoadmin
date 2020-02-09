package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"log"
	"strings"
)

type (
	ServiceStatus struct {
		Name string
		Status string
		RunningCount int
		PendingCount int
		DesiredCount int

	}
)

const (
	ECS_CLUSTER = "spotinst"
)

func ListServices() []string{
	sess := session.Must(session.NewSession())
	e := ecs.New(sess, aws.NewConfig())
	input := &ecs.ListServicesInput{
		Cluster: aws.String(ECS_CLUSTER),
	}
	out, err := e.ListServices(input)
	if err != nil {
		log.Println(err)
		return make([]string, 0)
	}
	result := make([]string, len(out.ServiceArns))
	for i, v := range out.ServiceArns {
		t := strings.Split(*v, "/")
		result[i] = t[len(t) - 1]
	}
	return result
}

func GetServiceStatus(name string) *ServiceStatus {
	sess := session.Must(session.NewSession())
	e := ecs.New(sess, aws.NewConfig())
	input := &ecs.DescribeServicesInput{
		Cluster: aws.String(ECS_CLUSTER),
		Services: []*string{aws.String(name)},
	}
	out, err := e.DescribeServices(input)
	if err != nil || len(out.Services) != 1{
		log.Println(err)
		return nil
	}

	t := out.Services[0]
	return &ServiceStatus{
		Name:         *t.ServiceName,
		Status:       *t.Status,
		RunningCount: int(*t.RunningCount),
		PendingCount: int(*t.PendingCount),
		DesiredCount: int(*t.DesiredCount),
	}
}

func ScaleService(name string, count int) bool {
	sess := session.Must(session.NewSession())
	e := ecs.New(sess, aws.NewConfig())
	describeInput := &ecs.DescribeServicesInput{
		Cluster: aws.String(ECS_CLUSTER),
		Services: []*string{aws.String(name)},
	}
	describeOut, err := e.DescribeServices(describeInput)
	if err != nil {
		log.Println(err)
		return false
	}
	if  len(describeOut.Services) == 0 || len(describeOut.Services) != 1 {
		log.Println("invalid service length")
		return false
	}
	target := describeOut.Services[0]
	if *target.ServiceName != name {
		log.Println("invalid service")
		return false
	}
	updateInput := &ecs.UpdateServiceInput{
		Cluster: aws.String(ECS_CLUSTER),
		Service: target.ServiceName,
		DesiredCount: aws.Int64(int64(count)),
		ForceNewDeployment: aws.Bool(true),
		HealthCheckGracePeriodSeconds: aws.Int64(int64(60)),
		DeploymentConfiguration: &ecs.DeploymentConfiguration{
			MaximumPercent:        aws.Int64(int64(200)),
			MinimumHealthyPercent: aws.Int64(int64(25)),
		},
	}
	_, err = e.UpdateService(updateInput)
	return err == nil
}