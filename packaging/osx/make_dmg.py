#!/usr/bin/python
#This script should be run from inside the packaging/osx folder.
#./make_osx.py <version> [<path-to-config>]
#
#Make sure you have the package 'icnsutils' installed.
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
                                      bin_files={'../../client/subshard.py': 'subshard',
                                                 '../../client/subshard_configurator.py': 'subshard_configurator',
                                                 'subshard-install-cert.sh': 'subshard-install-cert'},
                                      executable='subshard',
                                      data_files={'../../client/cr_theme': 'cr_theme',
                                                  '../../client/subshard_extension': 'subshard_extension'},
                                      icon='../deb/chromeball_google_chrome_poke_by_azerik92-d4c31vz.png')

    print osx_builder.package(version, config_path)
