# Client Installation

## Ubuntu

1. Install the client package. You can either construct this package from source by running the packaging script in packaging/deb, or downloading a prebuilt package from the github.
2. Install the TLS certificate for your client by invoking this on the command line: `sudo subshard-install-cert cert.pem` - where `cert.pem` is the path to the server certificate you downloaded while setting up the server.
3. Invoke subshard, either from the command line or your launcher.
4. On first run, you will need to enter the address and port of your server. Enter it like this: `https://address:port`. The port is `8080` by default.
5. Reinvoke subshard, and you should now be up and running!

## OSX

1. Open the .dmg image. You can either construct this image from source by running the packaging script in packaging/osx, or downloading a prebuilt package from the github.
2. Drag the subshard icon over to the application icon and release. Unmount the image.
3. Install the TLS certificate for your client by invoking this in terminal: `/Applications/Subshard.app/Contents/MacOS/bin/subshard-install-cert cert.pem` - where `cert.pem` is the path to the server certificate you downloaded while setting up the server. It will prompt you for your password, as it is installing the cert to your trust keychain.
4. Invoke subshard, either from the command line or your launcher.
5. On first run, you will need to enter the address and port of your server. Enter it like this: `https://address:port`. The port is `8080` by default.
6. Reinvoke subshard, and you should now be up and running!
