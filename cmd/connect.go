// Copyright © 2017 Jeremy Patton <jeremy@mutter.io>
//

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mutterio/aws-ssh/models"
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect to a known instance",
	Run:   execConnect,
}

var host string
var user string
var key string
var port string

func init() {
	RootCmd.AddCommand(connectCmd)
	connectCmd.PersistentFlags().StringVarP(&host, "server", "s", "", "server to connect to, can be IP or hostname")
	connectCmd.PersistentFlags().StringVarP(&user, "user", "u", "", "user to login with")
	connectCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "pem key name")
	connectCmd.PersistentFlags().StringVarP(&port, "port", "p", "22", "port to connect to")

}

func execConnect(c *cobra.Command, args []string) {
	instance := models.Instance{
		Host: host,
		Key:  key,
		User: user,
		Port: port,
	}
	fmt.Println("HOST::: ", host)
	Connect(instance)
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
	fmt.Println("Connecting to: ")
	fmt.Println("Name: ", inst.Name)
	fmt.Println("Private ip:", inst.PrivateIp)
	fmt.Println("Args:::", strings.Join(cmd.Args, " "))
	fmt.Println("##############################################")

	cmd.Run()
}
