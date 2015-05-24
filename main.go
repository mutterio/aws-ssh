package main // import "github.com/mutterio/aws-ssh"

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/ec2"
	"github.com/codegangsta/cli"
)

// Name is the exported name of this application.
const Name = "aws-ssh"

// Version is the current version of this application.
const Version = "0.0.2.dev"

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:  "search",
			Usage: "find a server and shell into it.",
			Action: func(c *cli.Context) {
				res := getInstances()
				instances := InstancesFromReservations(res, "")
				server := ""
				if len(c.Args()) > 0 {
					server = c.Args()[0]
				}
				findServer(server, instances)
			},
		},
		{
			Name:  "config",
			Usage: "generate aws ssh config",
			Action: func(c *cli.Context) {
				outDir := c.String("out")
				keyPath := c.String("keypath")
				res := getInstances()
				instances := InstancesFromReservations(res, keyPath)
				generateConfig(instances, outDir)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "out, o",
					Value: "",
					Usage: "file path to be written to",
				},
				cli.StringFlag{
					Name:  "keypath, kp",
					Value: "~/.ssh",
					Usage: "path for pem keys",
				},
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
					Host: c.String("host"),
					Key:  c.String("key"),
					User: c.String("user"),
					Port: c.String("port"),
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
	instance, err := selectInstance(server, instances)
	if err != nil {
		fmt.Println(err)
		return
	}
	shell(instance)
}

const hostTemplate = `
Host {{.Name}}
HostName {{.Host}}
User {{.User}}
EnableSSHKeysign yes
IdentityFile {{.KeyPath}}
`

func generateConfig(instances []Instance, outFile string) {
	t := template.Must(template.New("host").Parse(hostTemplate))
	var buf bytes.Buffer

	for _, inst := range instances {
		err := t.Execute(&buf, inst)
		if err != nil {
			fmt.Println("template err", err)
		}
	}

	if outFile != "" {

		ioutil.WriteFile(outFile, buf.Bytes(), 0600)
	} else {
		buf.WriteTo(os.Stdout)
	}

}

func selectInstance(server string, instances []Instance) (Instance, error) {
	matches := []Instance{}
	for _, instance := range instances {
		if strings.HasPrefix(instance.Name, server) {
			matches = append(matches, instance)
		}
	}
	if len(matches) == 1 {
		return Instance{}, errors.New("Server not found")
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	fmt.Println("Found ", len(matches), "matches in", len(instances), "instances")
	for pos, match := range matches {
		fmt.Println(pos, "  ", match.Name)
	}
	fmt.Print("Select vm: ")
	var input string
	fmt.Scanln(&input)
	fmt.Print(input)
	idx, err := strconv.ParseInt(input, 10, 8)
	if err != nil {
		log.Fatal(err)
	}
	return matches[idx], nil

}

func shell(inst Instance) {
	remoteServer := inst.Host
	if inst.User != "" {
		remoteServer = fmt.Sprintf("%v@%v", inst.User, remoteServer)
	}

	cmd := exec.Command("ssh")

	if inst.Key != "" {
		cmd.Args = append(cmd.Args, "-i")
		cmd.Args = append(cmd.Args, inst.KeyPath())
	}

	if inst.Port != "" {
		cmd.Args = append(cmd.Args, "-p")
		cmd.Args = append(cmd.Args, inst.Port)
	}

	cmd.Args = append(cmd.Args, remoteServer)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	fmt.Println("##############################################")
	fmt.Println("Connecting to: ", inst.Name)
	fmt.Println("Args:::", strings.Join(cmd.Args, " "))
	fmt.Println("##############################################")

	cmd.Run()
}

func getInstances() []*ec2.Reservation {
	fmt.Println("looking up instances.....")
	svc := ec2.New(&aws.Config{Region: "us-east-1"})

	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}
	return resp.Reservations

}
