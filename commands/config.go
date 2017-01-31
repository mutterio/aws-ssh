package commands

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/mutterio/aws-ssh/models"
	"github.com/urfave/cli"
)

var CmdConfig = cli.Command{
	Name:   "config",
	Usage:  "generate aws ssh config",
	Action: runConfig,
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
}

func runConfig(c *cli.Context) error {
	outDir := c.String("out")
	keyPath := c.String("keypath")
	region := c.GlobalString("region")
	res := models.GetInstances(region)
	instances := models.InstancesFromReservations(res, keyPath)
	generateConfig(instances, outDir)
	return nil
}

const hostTemplate = `
Host {{.Name}}
HostName {{.Host}}
User {{.User}}
EnableSSHKeysign yes
IdentityFile {{.KeyPath}}
`

func generateConfig(instances []models.Instance, outFile string) {
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
