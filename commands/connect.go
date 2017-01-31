package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mutterio/aws-ssh/models"
	"github.com/urfave/cli"
)

var CmdConnect = cli.Command{
	Name:   "connect",
	Usage:  "connect to a known instance",
	Action: runCmd,
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
}

func runCmd(c *cli.Context) error {
	fmt.Println(strings.Join(c.Args(), " "))
	instance := models.Instance{
		Host: c.String("host"),
		Key:  c.String("key"),
		User: c.String("user"),
		Port: c.String("port"),
	}
	fmt.Println("HOST::: ", c.String("host"))
	Connect(instance)
	return nil
}

func Connect(inst models.Instance) {
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
