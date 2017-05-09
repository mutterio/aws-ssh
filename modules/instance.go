package modules

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type Tag struct {
	Key   string
	Value string
}

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
	Tags        []Tag
}

//KeyPath get full key path for instance
func (i Instance) KeyPath(basePath string) string {
	if i.Key == "" {
		return ""
	}
	return fmt.Sprintf("%v/%v.pem", basePath, i.Key)
}

func (i Instance) GetKey(key string) string {
	for _, tag := range i.Tags {
		if tag.Key == key {
			return tag.Value
		}
	}
	return ""
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
	tags := []Tag{}
	for _, tag := range inst.Tags {
		tags = append(tags, Tag{Key: *tag.Key, Value: *tag.Value})
	}

	return Instance{
		Id:        *inst.InstanceId,
		Name:      name,
		Host:      host,
		Key:       key,
		State:     state,
		PrivateIp: privateIp,
		ImageId:   *inst.ImageId,
		Tags:      tags,
	}

}
