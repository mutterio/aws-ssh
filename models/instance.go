package models

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mitchellh/go-homedir"
)

//ec2 instance type
type Instance struct {
	Id          string
	User        string
	Host        string
	PrivateIp   string
	Key         string
	Name        string
	Port        string
	BaseKeyPath string
	State       string
	Type        string
}

func (i Instance) KeyPath() string {
	if i.Key == "" {
		return ""
	}
	basePath := i.BaseKeyPath
	if i.BaseKeyPath == "" {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		basePath = fmt.Sprintf("%v/.ssh", home)
	}

	return fmt.Sprintf("%v/%v.pem", basePath, i.Key)

}

func InstancesFromReservations(reservations []*ec2.Reservation, keyPath string) []Instance {
	instances := []Instance{}
	for _, res := range reservations {
		for _, inst := range res.Instances {
			name := "None"
			user := "ubuntu"
			key := ""
			host := ""
			privateIp := ""
			state := *inst.State.Name

			for _, keys := range inst.Tags {
				if *keys.Key == "Name" {
					name = *keys.Value
				}
				if *keys.Key == "User" {
					name = *keys.Value
				}
			}
			if inst.KeyName != nil {
				key = *inst.KeyName
			}
			if inst.PublicIpAddress != nil {
				host = *inst.PublicIpAddress
			}
			if inst.PrivateIpAddress != nil {
				privateIp = *inst.PrivateIpAddress
			}

			instances = append(instances, Instance{
				Id:          *inst.InstanceId,
				Name:        name,
				User:        user,
				Host:        host,
				Key:         key,
				BaseKeyPath: keyPath,
				State:       state,
				PrivateIp:   privateIp,
			})
		}
	}
	return instances
}

func GetInstances(region string) []*ec2.Reservation {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("looking up instances.....")
	cfgs := &aws.Config{Region: &region}
	svc := ec2.New(sess, cfgs)

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}
	return resp.Reservations

}
