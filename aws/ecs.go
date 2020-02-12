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
		Name         string
		Status       string
		RunningCount int
		PendingCount int
		DesiredCount int
		Tasks        []TaskInfo
	}
	TaskInfo struct {
		Name   string
		Image  string
		CPU    int
		Memory int
	}
)

const (
	ECS_CLUSTER = "spotinst"
)

func ListServices() []string {
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
		result[i] = t[len(t)-1]
	}
	return result
}

func GetServiceStatus(name string) *ServiceStatus {
	sess := session.Must(session.NewSession())
	e := ecs.New(sess, aws.NewConfig())
	serviceInput := &ecs.DescribeServicesInput{
		Cluster:  aws.String(ECS_CLUSTER),
		Services: []*string{aws.String(name)},
	}
	serviceOut, err := e.DescribeServices(serviceInput)
	if err != nil || len(serviceOut.Services) != 1 {
		log.Println(err)
		return nil
	}

	t := serviceOut.Services[0]

	tasks := make([]TaskInfo, 0)
	if len(t.Deployments) > 0 {
		taskDefinition := t.Deployments[0].TaskDefinition
		taskInput := &ecs.DescribeTaskDefinitionInput{
			TaskDefinition: taskDefinition,
		}
		taskOut, err := e.DescribeTaskDefinition(taskInput)
		if err != nil {
			log.Println(err)
			return nil
		}
		for _, v := range taskOut.TaskDefinition.ContainerDefinitions {
			tasks = append(tasks, TaskInfo{
				Name:   *v.Name,
				Image:  *v.Image,
				CPU:    int(*v.Cpu),
				Memory: int(*v.MemoryReservation),
			})
		}
	}

	return &ServiceStatus{
		Name:         *t.ServiceName,
		Status:       *t.Status,
		RunningCount: int(*t.RunningCount),
		PendingCount: int(*t.PendingCount),
		DesiredCount: int(*t.DesiredCount),
		Tasks:        tasks,
	}
}

func ScaleService(name string, count int) bool {
	sess := session.Must(session.NewSession())
	e := ecs.New(sess, aws.NewConfig())
	describeInput := &ecs.DescribeServicesInput{
		Cluster:  aws.String(ECS_CLUSTER),
		Services: []*string{aws.String(name)},
	}
	describeOut, err := e.DescribeServices(describeInput)
	if err != nil {
		log.Println(err)
		return false
	}
	if len(describeOut.Services) == 0 || len(describeOut.Services) != 1 {
		log.Println("invalid service length")
		return false
	}
	target := describeOut.Services[0]
	if *target.ServiceName != name {
		log.Println("invalid service")
		return false
	}
	updateInput := &ecs.UpdateServiceInput{
		Cluster:                       aws.String(ECS_CLUSTER),
		Service:                       target.ServiceName,
		DesiredCount:                  aws.Int64(int64(count)),
		ForceNewDeployment:            aws.Bool(true),
		HealthCheckGracePeriodSeconds: target.HealthCheckGracePeriodSeconds,
		DeploymentConfiguration:       target.DeploymentConfiguration,
	}
	_, err = e.UpdateService(updateInput)
	if err != nil {
		log.Println(err)
	}
	return err == nil
}
