// Copyright © 2017 Jeremy Patton <jeremy@mutter.io>
//

package cmd

import (
	"fmt"

	"path/filepath"

	"github.com/mutterio/aws-ssh/modules"
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
	keyDir, keyName := filepath.Split(key)
	instance := modules.Instance{
		Host: host,
		Key:  keyName,
		User: user,
		Port: port,
	}
	fmt.Println("HOST::: ", host)

	modules.Connect(instance, keyDir)
}
