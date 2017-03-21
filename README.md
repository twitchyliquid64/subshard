# Subshard <img src="https://github.com/twitchyliquid64/subshard/raw/master/packaging/deb/chromeball_google_chrome_poke_by_azerik92-d4c31vz.png" width="48">

Subshard is an attempt to make a better 'Tor browser', based on Chrome instead of Firefox, and using a client-server model instead of a 'thick client'.

## How it works

Subshard has two parts: The client script (which launches and controls Chrome for you) and `subshard-serv`, where Chrome tunnels all traffic.

When you launch subshard, a separate chrome instance starts which is bound to a subshard server. All traffic goes through your server, proxy style.

You configure which domains go to Tor and which URLs go to the open internet via the server (not through Tor). You can configure domains to be blocked, such as ad domains or websites which can identify you (prevent yourself from absent-mindedly logging into facebook).

## Why do you believe this is better than the Tor bundle?

_The Tor browser is an easy target for browser exploits. Such approaches are known and effective attack vectors by modern adversaries._

 * Firefox lacks a lot of generic exploit protections that other browsers implement.
 * The Tor browser cuts its updates from Extended Support Release, often getting important security fixes _months_ after it is patched. Chrome, by contrast prioritises these releases and has them out typically in 2 weeks.
 * There are comparitively few versions of the browser bundle, meaning exploit writers have an easier job targeting exploits for Tor users.
 * Subshard is an attempt to bring the best of Chromes security into the Tor ecosystem.

## Other features

 * Entirely separate chrome instance - no sharing of cookies / history / local storage / extensions.
 * Bright red color theme means it is always obvious which chrome you are in.
 * Serverside domain blacklists can help to prevent accidental browsing to certain sites, or prevent traffic hitting certain domains (eg: ad domains).
 * Automatically forward traffic to specific domains to another SOCKS/HTTP/HTTPS proxy. These rules are defined by regexes on the domain.
 * Multiple user support
 * We have an extension - `subshard guard` - that hacks in basic first-party isolation for Tor domains. This is the one feature of Firefox that I miss.

## Installation

### Server

Please see our instructions on [installing the server](https://github.com/twitchyliquid64/subshard/blob/master/server_install.MD).

### Client

__Windows:__ Unfortunately, I dont have packaging for Windows. Please [see here](https://github.com/twitchyliquid64/subshard/blob/master/windows_install.MD) for some basic instructions for manually getting it working. If you know windows, please help me get it packaged.

__OSX/Debian:__ Setting up a client on your machine should be as trivial as installing a package for your OS, doing two configuration steps for your first run, and then using it!
Make sure you have Chrome stable already installed, and please follow our [instructions](https://github.com/twitchyliquid64/subshard/blob/master/client_install.MD) for your platform.


### Troubleshooting

For issues in your Client: [Client Troubleshooting](https://github.com/twitchyliquid64/subshard/blob/master/troubleshooting.MD)

## References

 * [Hacking Subshard](https://github.com/twitchyliquid64/subshard/blob/master/hacking.MD) - A lay of the ~~land~~ source + design
 * [Installing the server](https://github.com/twitchyliquid64/subshard/blob/master/server_install.MD)

## I need your help!

 - [ ] (!) It needs to be packaged (to a zip, msi, something) for Windows. Help?
 - [ ] Make options page for Subshard Guard Chrome Extension work.
 - [ ] Make symlinks for our binaries in /usr/local/bin on OSX.
 - [ ] Some code cleanups around proxy initialization and Authentication.
 - [ ] Support for Client certificate authentication.
 - [ ] Automated testing for important security features, rather than me manually testing each release.
 - [ ] More security features!
