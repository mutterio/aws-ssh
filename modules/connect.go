package modules

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Connect(inst Instance, keyPath string) {
	remoteServer := inst.Host
	if inst.User != "" {
		remoteServer = fmt.Sprintf("%v@%v", inst.User, remoteServer)
	}

	cmd := exec.Command("ssh")

	if inst.Key != "" {
		cmd.Args = append(cmd.Args, "-i")
		cmd.Args = append(cmd.Args, inst.KeyPath(keyPath))
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
