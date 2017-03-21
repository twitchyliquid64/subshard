# Hacking

## Server Overview

[main.go](https://github.com/twitchyliquid64/subshard/blob/master/serv/main.go) is the entrypoint of the server. It ingests the config, constructs rules to handle specific URLs using a closure-like pattern ([proxy_init.go](https://github.com/twitchyliquid64/subshard/blob/master/serv/proxy_init.go)), opens a listener with specific TLS configuration ([listener_init.go](https://github.com/twitchyliquid64/subshard/blob/master/serv/listener_init.go)), and then starts serving. OS signals either terminate the server, or SIGHUP reopens the listener and reloads the configuration.

 * You can add new fields to the configuration file  [here](https://github.com/twitchyliquid64/subshard/blob/master/serv/config.go)

## Client Overview

Really, all we do is launch Chrome from the command line with a lot of flags. This is all implemented as a [python script](https://github.com/twitchyliquid64/subshard/blob/master/client/subshard.py).

Importantly, we:

 * Read the system subshard configuration (/etc/subshard/subshard.json), then read the users subshard configuration.
 * If this is the firstrun (~/.subshard doesn't exist) we instead launch the firstrun configuration utility [subshard_configurator.py](https://github.com/twitchyliquid64/subshard/blob/master/client/subshard_configurator.py)
 * Start chrome with specific args
  * User data dir - separate directory all of Chromes data. Completely separates it from your normal usage of Chrome.
  * Proxy settings - force Chrome to use the subshard server
  * `load_extension` - force Chrome to load our Theme and our extension, Subshard Guard.

## Building

Building the server is easy, you cd into `serv/`, set the gopath ``export GOPATH=`pwd` ``, and `go build`. `serv` should now exist. To make the deb for the server, you cd into `package/server_deb` and invoke `make_deb.py <version> <architecture>`.

Building the client packages is more complex, you cd into `packaging/{deb,dmg}` and run the corresponding make python script. This was built to be run in a debian environment (especially make_dmg.py), so you may be missing some commands. See the python script for details.
