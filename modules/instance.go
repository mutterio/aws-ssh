package modules

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

//Instance ec2 instance type
type Instance struct {
	Id          string
	User        string
	Host        string
	PrivateIp   string
	Key         string
	Name        string
	Port        string
	State       string
	Type        string
	ImageId     string
	FullKeyPath string
}

//KeyPath get full key path for instance
func (i Instance) KeyPath(basePath string) string {
	if i.Key == "" {
		return ""
	}
	return fmt.Sprintf("%v/%v.pem", basePath, i.Key)
}

func convertFromAwsType(inst *ec2.Instance) Instance {
	name := "None"
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
	return Instance{
		Id:        *inst.InstanceId,
		Name:      name,
		Host:      host,
		Key:       key,
		State:     state,
		PrivateIp: privateIp,
		ImageId:   *inst.ImageId,
	}
}
