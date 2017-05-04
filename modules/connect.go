package modules

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func getHost(inst Instance, user string) string {
	var host string
	if inst.Host == "" {
		host = inst.PrivateIp
	} else {
		host = inst.Host
	}

	if user != "" {
		host = fmt.Sprintf("%v@%v", user, host)
	} else if inst.User != "" && user != "" {
		host = fmt.Sprintf("%v@%v", inst.User, host)
	}

	return host

}

func expand(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

func Connect(inst Instance, user string, keyPath string) {
	remoteServer := getHost(inst, user)

	cmd := exec.Command("ssh")

	if inst.Key != "" {
		cmd.Args = append(cmd.Args, "-i")
		key, err := expand(inst.KeyPath(keyPath))
		if err != nil {
			panic(err)
		}
		cmd.Args = append(cmd.Args, key)
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
