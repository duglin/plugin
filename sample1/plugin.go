package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"plugin"
)

func main() {
	if err := plugin.Register(Start, Request); err != nil {
		fmt.Printf("Error registering: %q\n", err)
	}
}

// Start is executed when the plugin is first started in case we need to
// do some processing even before the main talks to us
func Start(c *plugin.Main) {
	// In this case we'll just ask the main for the DocerHost value
	buf, err := c.Call("GetDockerHost", "")
	if err != nil {
		fmt.Printf("Error tryign to get DockerHost value from main: %q\n", err)
	} else {
		fmt.Printf("test.Start.Dockerhost: %s\n", string(buf))
	}
}

var inCount int

// Request is the func that is called to process an incoming message.
func Request(c *plugin.Main, cmd string, buf []byte) ([]byte, error) {
	switch cmd {
	case plugin.CmdGetMetadata:
		// The fields in this metadata (map) will be plugin-driver
		// specific.  Each type of plugin will require a different set
		// of metadata. In this case we'll use some that the CLI might need
		md := map[string]string{
			"Cmd":  "networks",
			"Help": "me Rhonda",
		}
		return json.Marshal(md)

	// Note for this sample we'll just use strcmp but we could also
	// convert this to a typed function call if we wanted to add a
	// plugin-specific proxy layer - like what Docker daemon does
	case "Echo":
		if inCount == 0 {
			go doRepeat(c, buf)
		}
		inCount++
		return []byte(fmt.Sprintf("Response: %s", buf)), nil

	default:
		return nil, fmt.Errorf("No such cmd in plugin: %s", cmd)
	}
}

func doRepeat(c *plugin.Main, buf []byte) {
	done := make(chan int)
	count := 0
	max, _ := strconv.Atoi(strings.SplitN(string(buf), "/", 2)[1])
	fmt.Printf("Sending %d Repeats from plugin to main\n", max)
	for i := 1; i <= max; i++ {
		go func(i int) {
			req := fmt.Sprintf("Repeat me - %d/%s", i, max)
			res, err := c.Call("Repeat", req)
			if err != nil {
				fmt.Printf("Plugin Err %d: %q\n", i, err)
			} else if res != "Response: "+req {
				fmt.Printf("%d plugin: %s\n", i, res)
				os.Exit(0)
			} else {
				count++
				if count%(max/10) == 0 {
					fmt.Print(".")
				}
				if count == max {
					done <- 1
				}
			}
		}(i)
	}
	<-done
	c.Notify("Done", "")
}
