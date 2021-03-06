package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"plugin"
)

func main() {
	if err := plugin.Register(nil, Request); err != nil {
		fmt.Printf("Error registering: %q\n", err)
	}
}

func Request(c *plugin.Main, cmd string, buf []byte) ([]byte, error) {
	switch cmd {
	case plugin.CmdGetMetadata:
		md := map[string]string{
			"Cmd":         "compose",
			"Description": "Docker's 'compose' command",
			"Help":        "See 'docker compose --help' for more",
		}
		return json.Marshal(md)

	case "run":
		args := []string{}
		err := json.Unmarshal(buf, &args)
		if err != nil {
			return nil, err
		}

		cmd := exec.Command("docker-compose", args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()

		if err != nil {
			// Just blank out the error message so its not printed
			// by the Docker cli.  This assumes that the docker-compose
			// exec would have printed stuff to the screen already.
			err = fmt.Errorf("")
		}

		return nil, err

	default:
		return nil, fmt.Errorf("No such cmd in compose: %s", cmd)
	}
}
