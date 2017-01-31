// Copyright © 2017 Jeremy Patton <jeremy@mutter.io>
//

package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mutterio/aws-ssh/models"
	"github.com/olekukonko/tablewriter"
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
	searchCmd.PersistentFlags().StringVarP(&region, "region", "r", "us-east-1", "aws region to use")

}

func execSearch(c *cobra.Command, args []string) {
	res := models.GetInstances(region)
	instances := models.InstancesFromReservations(res, "")
	server := ""
	if len(args) > 0 {
		server = args[0]
	}
	findServer(server, instances)
}

func findServer(server string, instances []models.Instance) {
	instance, err := selectInstance(server, instances)
	if err != nil {
		fmt.Println(err)
		return
	}
	Connect(instance)
}

func selectInstance(server string, instances []models.Instance) (models.Instance, error) {
	matches := []models.Instance{}
	for _, instance := range instances {
		if strings.HasPrefix(instance.Name, server) {
			matches = append(matches, instance)
		}
	}
	if len(matches) == 0 {
		return models.Instance{}, errors.New("Server not found")
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	fmt.Println("Found ", len(matches), "matches in", len(instances), "instances")
	writeInstances(matches)
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

func writeInstances(matches []models.Instance) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Num", "Id", "State", "Public", "Private"})

	for pos, match := range matches {
		table.Append([]string{strconv.Itoa(pos), match.Id, match.State, match.Host, match.PrivateIp})
		// fmt.Println(pos, "  ", match.Name, " ", match.State, " ", match.Host)
	}
	table.Render()
}
