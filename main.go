package main // import "github.com/mutterio/aws-ssh"

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ec2"
	"github.com/codegangsta/cli"
	"github.com/mitchellh/go-homedir"
)

// Name is the exported name of this application.
const Name = "aws-ssh"

// Version is the current version of this application.
const Version = "0.0.1.dev"

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "search",
			Usage: "find a server and shell into it.",
			Action: func(c *cli.Context) {
				res := getInstances()
				instances := parseInstances(res)
				server := ""
				if len(c.Args()) > 0 {
					server = c.Args()[0]
				}
				findServer(server, instances)
			},
		},
		{
			Name:  "connect",
			Usage: "connect to a known instance",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host, hn",
					Value: "localhost",
					Usage: "server to connect to, can be IP or hostname",
				},
				cli.StringFlag{
					Name:  "user, u",
					Value: "",
					Usage: "user to login with",
				},
				cli.StringFlag{
					Name:  "key, k",
					Value: "",
					Usage: "pem key name",
				},
				cli.StringFlag{
					Name:  "port, p",
					Value: "22",
					Usage: "port to connect to",
				},
			},
			Action: func(c *cli.Context) {
				fmt.Println(strings.Join(c.Args(), " "))
				instance := Instance{
					host: c.String("host"),
					key:  c.String("key"),
					user: c.String("user"),
					port: c.String("port"),
				}
				fmt.Println("HOST::: ", c.String("host"))
				shell(instance)
			},
		},
	}
	app.Name = "aws-ssh"
	app.Usage = "Shell into a server!"

	app.Run(os.Args)
}

func findServer(server string, instances []Instance) {
	instance := selectInstance(server, instances)
	shell(instance)
}

//ec2 instance type
type Instance struct {
	user string
	host string
	key  string
	name string
	port string
}

func parseInstances(reservations []*ec2.Reservation) []Instance {
	instances := []Instance{}
	for _, res := range reservations {
		for _, inst := range res.Instances {
			name := "None"
			user := "ubuntu"
			key := ""
			host := ""
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
			if inst.PublicIPAddress != nil {
				host = *inst.PublicIPAddress
			}

			instances = append(instances, Instance{
				name: name,
				user: user,
				host: host,
				key:  key,
			})
		}
	}
	return instances
}
func selectInstance(server string, instances []Instance) Instance {
	matches := []Instance{}
	for _, instance := range instances {
		if strings.HasPrefix(instance.name, server) {
			matches = append(matches, instance)
		}
	}
	if len(matches) == 1 {
		return matches[0]
	}
	fmt.Println("Found ", len(matches), "matches in", len(instances), "instances")
	for pos, match := range matches {
		fmt.Println(pos, "  ", match.name)
	}
	fmt.Print("Select vm: ")
	var input string
	fmt.Scanln(&input)
	fmt.Print(input)
	idx, err := strconv.ParseInt(input, 10, 8)
	if err != nil {
		log.Fatal(err)
	}
	return matches[idx]

}

func shell(inst Instance) {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	remoteServer := inst.host
	if inst.user != "" {
		remoteServer = fmt.Sprintf("%v@%v", inst.user, remoteServer)
	}

	var keyPath string
	cmd := exec.Command("ssh")

	if inst.key != "" {
		keyPath = fmt.Sprintf("%v/.ssh/%v.pem", home, inst.key)
		cmd.Args = append(cmd.Args, "-i")
		cmd.Args = append(cmd.Args, keyPath)
	}

	if inst.port != "22" {
		cmd.Args = append(cmd.Args, "-p")
		cmd.Args = append(cmd.Args, inst.port)
	}

	cmd.Args = append(cmd.Args, remoteServer)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	fmt.Println("##############################################")
	fmt.Println("Connecting to: ", inst.name)
	fmt.Println("Args:::", strings.Join(cmd.Args, " "))
	fmt.Println("##############################################")

	cmd.Run()
}

func getCommand(inst *ec2.Instance) (string, string) {
	key := *inst.KeyName
	instanceUser := "ubuntu"
	ip := *inst.PublicIPAddress
	for _, keys := range inst.Tags {
		if *keys.Key == "User" {
			instanceUser = url.QueryEscape(*keys.Value)
		}
	}
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	home := usr.HomeDir
	keyPath := fmt.Sprintf("%v/.ssh/%v.pem", home, key)
	remote := fmt.Sprintf("%v@%v", instanceUser, ip)

	return keyPath, remote
}

func getInstances() []*ec2.Reservation {
	fmt.Println("looking up instances.....")
	svc := ec2.New(&aws.Config{Region: "us-east-1"})

	// Call the DescribeInstances Operation
	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}
	return resp.Reservations

}
