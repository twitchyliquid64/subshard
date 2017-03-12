# subshard

Subshard is a well packaged, secure-by-default web proxy. Deploy the server component on a remote box/VPS/EC2. Install the client component on your computer. Browse the web by
proxy, even setup forwarders on the server to allow you to browse darknets without messy config on your client.

## Installation

### Server

_Requires Ubuntu. Can be trivally made to work on other UNIX platforms but you will have to setup the certs/config files yourself._

1. Setup a DNS address for your server if it doesnt have one already. In other words, make sure you can access it by name (proxy.example.com) rather than IP address.
2. Install the server package for your system - eg: `sudo dpkg -i server-amd64-0.1.0.deb`. You can either construct this package from source by building the server and running
the packaging script, or downloading a prebuilt package from the github.
3. Follow the prompts to setup a default configuration. Specifically, create a username and password for first use, and put in the DNS address from step 1.
4. Copy `/etc/subshard/cert.pem` to your machine. Each client will need a copy of this file.
5. Everything should now be setup! You can now run `sudo service subshard-serv start` to bring it up, and `sudo service subshard-serv stop` to bring it down.

#### Server administration

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

You can do more on the config file directly, which is located at `/etc/subshard/subshard-serv.json`. I do not recommend this however, its possible to have an insecure configuration.

### Client

Setting up a client on your machine should be as trivial as installing a package for your OS, doing two configuration steps for your first run, and then using it!
Make sure you have Chrome already installed.

#### Ubuntu

1. Install the client package. You can either construct this package from source by running the packaging script in packaging/deb, or downloading a prebuilt package from the github.
2. Install the TLS certificate for your client by invoking this on the command line: `sudo subshard-install-cert cert.pem` - where `cert.pem` is the path to the server certificate you downloaded while setting up the server.
3. Invoke subshard, either from the command line or your launcher.
4. On first run, you will need to enter the address and port of your server. Enter it like this: `http://address:port`. The port is `8080` by default.
5. Reinvoke subshard, and you should now be up and running!

#### OSX

1. Open the .dmg image. You can either construct this image from source by running the packaging script in packaging/osx, or downloading a prebuilt package from the github.
2. Drag the subshard icon over to the application icon and release. Unmount the image.
3. Install the TLS certificate for your client by invoking this in terminal: `/Applications/Subshard.app/Contents/MacOS/bin/subshard-install-cert cert.pem` - where `cert.pem` is the path to the server certificate you downloaded while setting up the server. It will prompt you for your password, as it is installing the cert to your trust keychain.
4. Invoke subshard, either from the command line or your launcher.
5. On first run, you will need to enter the address and port of your server. Enter it like this: `http://address:port`. The port is `8080` by default.
6. Reinvoke subshard, and you should now be up and running!

#### Troubleshooting

##### Chrome error: ERR_PROXY_CERTIFICATE_INVALID

You need to install the servers certificate into the trust root of your computer, because otherwise Chrome won't trust your server. This can be done by invoking `subshard-install-cert <certificate-file>` from the command line.

##### I need to change the server address

1. Invoke `subshard_configurator` from the command line.
2. Input the key 'proxy_addr' and press enter.
3. Input the new server address in the format 'https://address:port' and press enter.
3. Control-C to exit the configurator.

## Development TODO

 - [ ] Monitor health of a forwarder
 - [ ] Verification endpoint for Subshard Guard.
 - [ ] Make options page for Subshard Guard work.
 - [ ] Use a shell script to make symlinks in /usr/local/bin on OSX
