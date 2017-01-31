package main

import (
	"os"

	"github.com/mutterio/aws-ssh/commands"
	"github.com/urfave/cli"
)

// Name is the exported name of this application.
const Name = "aws-ssh"

// Version is the current version of this application.
const Version = "0.0.6"

func main() {
	app := cli.NewApp()
	app.Name = "aws-ssh"
	app.Usage = "Shell into a server!"
	app.Commands = []cli.Command{
		commands.CmdSearch,
		commands.CmdConfig,
		commands.CmdConnect,
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "region, r",
			Value:  "us-east-1",
			EnvVar: "AWS_DEFAULT_REGION",
		},
	}
	app.Version = Version
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, V",
		Usage: "print only the version",
	}
	app.Run(os.Args)
}
