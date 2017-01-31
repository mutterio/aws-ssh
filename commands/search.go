package commands

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mutterio/aws-ssh/models"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var CmdSearch = cli.Command{
	Name:   "search",
	Usage:  "find a server and shell into it.",
	Action: runSearch,
}

func runSearch(c *cli.Context) error {
	region := c.GlobalString("region")
	res := models.GetInstances(region)
	instances := models.InstancesFromReservations(res, "")
	server := ""
	if len(c.Args()) > 0 {
		server = c.Args()[0]
	}
	findServer(server, instances)
	return nil
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
	if len(matches) == 1 {
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
