// Copyright © 2017 Jeremy Patton <jeremy@mutter.io>
//

package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/mutterio/aws-ssh/modules"
	"github.com/spf13/cobra"
)

var region string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long: `queries aws for instances and returns an ordered list allowing
  you to select an instance by number`,
	Run: execSearch,
}

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.PersistentFlags().StringVarP(&user, "user", "u", "user", "user to login with")
	searchCmd.PersistentFlags().StringVarP(&region, "region", "r", "us-east-1", "aws region to use")
	searchCmd.PersistentFlags().StringVarP(&keyPath, "keypath", "k", "~/.ssh", "path for pem keys")

}

func execSearch(c *cobra.Command, args []string) {
	fmt.Println("region>>>", region)
	instances := modules.GetInstances(region)

	server := ""
	if len(args) > 0 {
		server = args[0]
	}
	findServer(server, instances)
}

func findServer(server string, instances modules.Instances) {
	instance, err := selectInstance(server, instances)
	if err != nil {
		fmt.Println(err)
		return
	}
	modules.Connect(instance, keyPath)
}

func selectInstance(server string, instances modules.Instances) (modules.Instance, error) {
	matches := instances.FilterByName(server)
	if len(matches) == 0 {
		return modules.Instance{}, errors.New("No instances Found")
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(server) > 0 {
		fmt.Println("Found ", len(matches), "matches in", len(instances), "instances")
	}
	matches.CreateTable(os.Stdout)
	fmt.Print("Select vm Num : ")
	var input string
	fmt.Scanln(&input)
	fmt.Print(input)
	idx, err := strconv.ParseInt(input, 10, 8)
	if err != nil {
		log.Fatal(err)
	}
	return matches[idx], nil

}
