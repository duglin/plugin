package main

import (
	"fmt"
	"os"

	"plugin"
)

var echoDone = make(chan int)
var repeatDone = make(chan int)

func main() {
	// Define a new plugin
	// Normally the 'main' would discover all of the registered
	// plugins (e.g. look in some config dir) and then call
	// NewPlug() for each one
	pi := plugin.NewPlugin("./plugin", Request)

	// It may not actually start each one until needed though
	if err := pi.Start(); err != nil {
		fmt.Printf("Err starting plugin: %q\n", err)
	}

	// Show some Config data that we got from the plugin
	fmt.Printf("cli.pi.Cmd: %s\n", pi.Config["Cmd"])
	fmt.Printf("cli.pi.Help: %s\n", pi.Config["Help"])

	// Do some concurrency testing
	max := 5000
	echoCount := 0
	fmt.Printf("\nSending %d Echos from main to plugin\n", max)
	for i := 1; i <= max; i++ {
		go func(i int) {
			req := fmt.Sprintf("Hi - %d/%d", i, max)
			res, err := pi.Call("Echo", req)
			if err != nil {
				fmt.Printf("Err %d: %q\n", i, err)
			} else if res != "Response: "+req {
				fmt.Printf("%d cli: %s\n", i, res)
				os.Exit(0)
			} else {
				echoCount++
				if echoCount%(max/10) == 0 {
					fmt.Print(".")
				}
				if echoCount == max {
					echoDone <- 1
				}
			}
		}(i)
	}

	// Give all threads a chance to complete
	<-echoDone
	<-repeatDone

	pi.Stop()
	fmt.Printf("\nDone\n")
}

func Request(p *plugin.Plugin, cmd string, buf []byte) ([]byte, error) {
	// In this case we just do strcmp's and do the work. But we could
	// also provide a layer that converts this to a normal
	// func() call with typed parameters
	if cmd == "GetDockerHost" {
		return []byte("127.0.0.1:2375"), nil
	}
	if cmd == "Repeat" {
		return []byte(fmt.Sprintf("Response: %s", buf)), nil
	}
	if cmd == "Done" {
		repeatDone <- 1
		return nil, nil
	}
	return nil, fmt.Errorf("No such comamnd: %s", cmd)
}
