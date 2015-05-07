This sample just shows the basics of how to use the plugin infrastructure.

It:
- has the main start the test plugin (in this case it simulates a CLI exe)
- the plugin will immediate ask for info from the main
- the main will kick off a series of 5000 parallel requests to the plugin
- the plugin, after getting the first of those parallel requests, will
start its own set of parallel requests (same #) back to the main

To play with it use 'mk'. This will compile and then run 'cli' which
is the sample 'main'. 
