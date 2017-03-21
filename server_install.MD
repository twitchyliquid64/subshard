# Server installation

_Requires Ubuntu. Can be trivally made to work on other UNIX platforms but you will have to setup the certs/config files yourself._

1. Setup a DNS address for your server if it doesnt have one already. In other words, make sure you can access it by name (proxy.example.com) rather than IP address.
2. Install the server package for your system - eg: `sudo dpkg -i server-amd64-0.1.0.deb`. You can either construct this package from source by building the server and running
the packaging script, or downloading a prebuilt package from the github.
3. Follow the prompts to setup a default configuration. Specifically, create a username and password for first use, and put in the DNS address from step 1.
4. Copy `/etc/subshard/cert.pem` to your machine. Each client will need a copy of this file.
5. Everything should now be setup! You can now run ~~`sudo service subshard-serv start` to bring it up, and `sudo service subshard-serv stop` to bring it down.~~ Manually run `subshard-serv` inside a screen session or in the background. I need to fix the sysv wrapper.

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
subshard-admin setlistener :5000 #Setup server on port 5000
```

You can do more on the config file directly, which is located at `/etc/subshard/subshard-serv.json`. I do not recommend this however, its possible to have an insecure configuration. All the fields in the configuration file are documented [here](https://github.com/twitchyliquid64/subshard/blob/master/serv/config.go).
