// Copyright © 2017 Jeremy Patton <jeremy@mutter.io>
//

package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

const Version = "0.0.8"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints the current version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("aws-ssh v%s %s/%s \n", Version, runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)

}
