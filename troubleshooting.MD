
## Chrome error: ERR_PROXY_CERTIFICATE_INVALID

You need to install the servers certificate into the trust root of your computer, because otherwise Chrome won't trust your server. This can be done by invoking `subshard-install-cert <certificate-file>` from the command line.

## Chrome error: ERR_PROXY_CONNECTION_FAILED

Your server is not running, is not accessible, or the server address is incorrect. Confirm the server is running, and check the server address as per instructions below.

## I need to change the server address on the client

1. Invoke `subshard_configurator` from the command line.
2. Input the key 'proxy_addr' and press enter.
3. Input the new server address in the format 'https://address:port' and press enter.
3. Control-C to exit the configurator.
