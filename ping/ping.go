package main

import (
	"encoding/json"
	"fmt"

	"plugin"
)

func main() {
	if err := plugin.Register(Start, Request); err != nil {
		fmt.Printf("Error registering: %q\n", err)
	}
}

func Start(c *plugin.Main) {
}

func Request(c *plugin.Main, cmd string, buf []byte) ([]byte, error) {
	switch cmd {
	case plugin.CmdGetMetadata:
		md := map[string]string{
			"Cmd":         "ping",
			"Description": "A ping command",
			"Help":        "Cool ping extension\nWhat it does do... I dunno",
		}
		return json.Marshal(md)

	case "run":
		host, err := c.Call("GetDockerHost", "")
		if err != nil {
			return nil, err
		}

		fmt.Printf("Activating ping...\n")
		fmt.Printf("DockerHost: %s\n", host)

		type callArgs struct {
			Method string
			Path   string
			Data   []byte
		}

		args := callArgs{
			Method: "GET",
			Path:   "/info",
			Data:   nil,
		}

		buf, err := json.Marshal(args)
		if err != nil {
			return nil, err
		}

		buf, err = c.CallBytes("CallDaemon", buf)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Daemon Info: %s\n", string(buf))

		return nil, nil

	default:
		return nil, fmt.Errorf("No such cmd in ping: %s", cmd)
	}
}
