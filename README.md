# Plugin

# Basic Concept
One exe (main) will start a plugin (2nd exe) and they'll communicate with
each other via pipes.

The communication is bi-directional (each side can initiate an exchange)
and multiple requests can happen at the same time.

It supports one-way and request-response exchanges.

The plugin is told what file-descriptors to use but if not then it just
uses stdin/out.

Right now the plugin is started as a new exe, but we could make it 
configurable so that its run within a Docker container.

Plugin discovery is not part of this yet since I think that might be
'main'-specific. But I could provide a default search location/method if wanted.

See the code and samples for more details - until I get around to
creating better docs :-)

# Playing with it...
See sample1

# How it works

## Transport Layer

The actual transport of messages is done by converting a outgoing messages
into a 'chunk'.  Each chunk consists of the following pieces of information:

* `threadID` - a unique ID that is used to correlate requests and responses.
* `direction` - indicated whether this chunk is a request, a response or a
  oneway message. A 'oneway' message is one that is like a notification,
  it is sent and no response message is expected. It is an integer.
* `cmd` - a string representing the command/operation that is to be invoked.
  In the case of a response message, the `cmd` will be the same as the 
  request's `cmd` field.
* `buffer` - a `[]byte` that holds the actual data of the message.

In many ways a chunk is very like an HTTP request. The `cmd` is very similar
to the `URL` and the `buffer` is similar to the body.  The other data is just
metadata that might appear in HTTP headers.

Each chunk is sent as a series of byte as shown below:
```
+-----------------+
| threadID + \0   |
+-----------------+
| direction + \0  |
+-----------------+
| len of cmd + \0 |
+-----------------+
| cmd             |
+-----------------+
| buffer          |
+-----------------+
```

`threadID`, `direction`, len of `cmd` and `cmd` are all strings followed
by a zero byte. This is done so that we can support any length for those
fields - meaning we parse the chars until we get a \0 and then interpet
the data read in.

### Sending and Receiving Chunks

When a message (chunk) is sent, if the outgoing chunk is a `request`
then code will block waiting for a response.

To help there are a series of utility functions available. Ones that
start with `Call` are for request/response message exchanges.  Ones that
start with `Notify` are for one-ways.  Each variant supports specifying
the data as either strings or []bytes.  If necessary, you can also
pass in a pointer to a `Chunk` itself, otherwise one will be created 
for you.

### Processing incoming messages

Upon receiving an incoming message the processor will invoke a call-up
function that was provided when things were setup.  This function is
expected to do whatever processing is necessary on the incoming chunk.
If an error is returned, and the incoming Chunk was a request message,
then the error is sent back to the sending end and an error is generated
on that side.

## Plugin layer

The plugin layer uses the `net`working layer to communicate between a
main executable and a plugin.  The plugin layer is responsible for starting
the plugin executable, establishing the conection and then handing off
any incoming chunk to the executables for processing.

### Plugin struct

The `Plugin` struct is provided to the main executable and can be used to
access information about the plugin as well as to communicate with it
by simply calling the appropriate `Call` or `Notify` function.

Once a `NewPlugin` is called to create a new plugin, the `Start` function
can be used to start the plugin and establish the connection. There is also
a `Stop` function that can be called to kill the plugin executable.

### Main struct

Like the `Plugin` struct, the `Main` struct is provided to the plugin
executable as a way to talk back to the main executable.  Once started
the plugin is expected to call `Register` to let the plugin infrastructure
know about its call-back functions.

### Call-backs

There are two types of call-backs.  Both the main and plugin exectables
can specify a call-back function to be invoked upon receipt of an incoming
message. The plugin executable can also specify an additional call-back
function that is invoked at the end of Register() processing. This allows
the plugin to take any action that needs to happen prior to any incoming
message.
