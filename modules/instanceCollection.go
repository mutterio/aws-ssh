package modules

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/olekukonko/tablewriter"
)

//Instances Instance Collection type
type Instances []Instance

//GetInstances find by region
func GetInstances(region string) Instances {
	res := getReservations(region)
	instances := instancesFromReservations(res)
	instances.fillImageDetails(region)

	return instances
}

func mapDescriptionToUser(description *string) string {
	desc := ""
	if description != nil {
		desc = *description
	}
	if strings.Contains(desc, "Ubuntu") {
		return "ubuntu"
	}
	if strings.Contains(desc, "Centos") {
		return "centos"
	}
	return "ec2-user"
}

func (instances Instances) fillImageDetails(region string) {
	imageIds := instances.CollectString(func(inst Instance) string {
		return inst.ImageId
	})
	uniqueIds := RemoveDuplicatesUnordered(imageIds)
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	cfgs := &aws.Config{Region: &region}
	svc := ec2.New(sess, cfgs)
	input := &ec2.DescribeImagesInput{ImageIds: uniqueIds}
	result, err := svc.DescribeImages(input)
	for _, image := range result.Images {
		for idx, inst := range instances {
			if inst.ImageId == *image.ImageId {
				instances[idx].User = mapDescriptionToUser(image.Description)
			}
		}
	}
}

func getReservations(region string) []*ec2.Reservation {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("Fetching instances.....")
	cfgs := &aws.Config{Region: &region}
	svc := ec2.New(sess, cfgs)

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}
	return resp.Reservations

}

func instancesFromReservations(reservations []*ec2.Reservation) Instances {
	instances := []Instance{}
	for _, res := range reservations {
		for _, inst := range res.Instances {
			instances = append(instances, convertFromAwsType(inst))
		}
	}
	return instances
}

//CollectString collects
func (instances Instances) CollectString(f func(Instance) string) []string {
	result := make([]string, len(instances))
	for i, item := range instances {
		result[i] = f(item)
	}
	return result
}

//FilterByName filters collection by instance name
func (instances Instances) FilterByName(name string) Instances {
	return instances.Filter(func(i Instance) bool {
		return strings.Contains(i.Name, name)
	})
}

//CreateTable creates a formated table of instances
func (instances Instances) CreateTable(writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Num", "Id", "State", "Public", "Private", "env", "role"})

	for pos, inst := range instances {
		table.Append([]string{strconv.Itoa(pos),
			inst.Id,
			inst.State,
			inst.Host,
			inst.PrivateIp,
			inst.GetKey("env"),
			inst.GetKey("role"),
		})
		// fmt.Println(pos, "  ", inst.Name, " ", inst.State, " ", inst.Host)
	}
	table.Render()
}

func (instances Instances) Filter(f func(Instance) bool) Instances {
	var vsf Instances
	for _, v := range instances {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
