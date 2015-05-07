Docker `ping` extension

This plugin is meant to be used by a Docker CLI. Just compile and place
the executabel in the ~/.docker/plugins/ directory.

This plugin will do a couple of things:
- ask the Docker CLI for the DockerHost value, and print it
- invoke the /info API on the daemon and print the results

This demonstrates how from within the plugin you can call back into the
CLI to either retrieve information or to actually have it talk to the
daemon for you.  Which is better than calling the daemon from the plugin
itself because the cli will handle the security and extra bits, like adding
the HTTP Headers from the config file.
