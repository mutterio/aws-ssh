// Copyright © 2017 Jeremy Patton <jeremy@mutter.io>
//

package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/mutterio/aws-ssh/modules"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "generate aws ssh config",
	Run:   execConfig,
}

var outDir string

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.PersistentFlags().StringVarP(&outDir, "out", "o", "", "file path to be written to")
}

func execConfig(c *cobra.Command, args []string) {
	instances := modules.GetInstances(region)
	generateConfig(instances, outDir)
}

const hostTemplate = `
Host {{.Name}}
HostName {{.Host}}
User {{.User}}
EnableSSHKeysign yes
IdentityFile {{.FullKeyPath}}
`

func generateConfig(instances []modules.Instance, outFile string) {
	t := template.Must(template.New("host").Parse(hostTemplate))
	var buf bytes.Buffer

	for _, inst := range instances {
		inst.FullKeyPath = inst.KeyPath(keyPath)
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
