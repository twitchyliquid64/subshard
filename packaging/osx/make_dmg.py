#!/usr/bin/python
#This script should be run from inside the packaging/osx folder.
#./make_osx.py <version> [<path-to-config>]
import sys
sys.path.append('..')
import packager

if __name__ == '__main__':
    version = sys.argv[1]
    config_path = None
    if len(sys.argv) > 2:
        config_path = sys.argv[2]

    osx_builder = packager.OSXPackage('Subshard',
                                      maintainer='Twitchyliquid64 <twitchyliquid64@ciphersink.net>',
                                      description='Subshard is an isolated chrome instance that tunnels all its traffic through a proxy.',
                                      bin_files={'../../client/subshard.py': 'subshard'},
                                      data_files={'../../serv/web': 'web',
                                                  '../deb/chromeball_google_chrome_poke_by_azerik92-d4c31vz.png': 'Subshard'},
                                      icon='/Applications/Subshard.app/Contents/Resources/Subshard')

    print osx_builder.package(version, config_path)
