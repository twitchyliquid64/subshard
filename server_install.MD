# Server installation

_Requires Ubuntu. Can be trivally made to work on other UNIX platforms but you will have to setup the certs/config files yourself._

1. Setup a DNS address for your server if it doesnt have one already. In other words, make sure you can access it by name (proxy.example.com) rather than IP address.
2. Install the server package for your system - eg: `sudo dpkg -i server-amd64-0.1.0.deb`. You can either construct this package from source by building the server and running
the packaging script, or downloading a prebuilt package from the github.
3. Follow the prompts to setup a default configuration. Specifically, create a username and password for first use, and put in the DNS address from step 1.
4. Copy `/etc/subshard/cert.pem` to your machine. Each client will need a copy of this file.
5. The proxy should now function! You can now run ~~`sudo service subshard-serv start` to bring it up, and `sudo service subshard-serv stop` to bring it down.~~ Manually run `subshard-serv` inside a screen session or in the background. I need to fix the sysv wrapper.

## To get Tor working

The steps above will get you a Proxy. To funnel everything in Tor, do the following:

1. Install Tor: `sudo apt-get install tor`
2. Launch Tor if it hasnt already: `sudo service start tor`
3. Now, we tell subshard to forward traffic matching a specific regex to Tor. Run `subshard-admin forward ".*" Tor localhost:9050 SOCKS`. Dont forget quotes around the regex!
 * This tells Subshard to add a new forwarding rule called 'Tor', matching all hosts that match `.*`, forwarding them to a SOCKS proxy at `localhost:9050`. This is the default listening address of tor.
 * If you wanted, you can instead only route `.onion` addresses through Tor and simply proxy the rest of the internet. This is not secure, so make an evaluation based on what you are doing.

TODO: Describe how to setup the health checker feature.

## Server administration

The server component ships with a configuration utility that you can use to add forwarders, blacklist domains, users, etc. Changes should be in effect for all new
subshard sessions.

Some examples:

```shell
subshard-admin --help
subshard-admin adduser user2
subshard-admin deluser user2
subshard-admin setpw user2
subshard-admin blacklist host-expression .*facebook.com
subshard-admin unblacklist host-expression .*facebook.com
subshard-admin setlistener :5000 #Set server listener on port 5000 instead

#Setup a forwarder, forwarding all URLs to a SOCKS proxy on localhost:9050.
subshard-admin forward ".*" MyOtherProxy localhost:9050 SOCKS
 #Routes all URLs matching the pattern .*.internal, to MyProxy. Assumes MyProxy forwarder already exists.
subshard-admin forward ".*.internal" MyOtherProxy
#Remotes the routing rule matching '.*' to the MyOtherProxy forwarder.
subshard-admin unforward ".*" MyOtherProxy
```

You can do more on the config file directly, which is located at `/etc/subshard/subshard-serv.json`. I do not recommend this however, its possible to have an insecure configuration. All the fields in the configuration file are documented [here](https://github.com/twitchyliquid64/subshard/blob/master/serv/config.go).
