#!/usr/bin/python
#This script should be run from inside the packaging/deb folder.
#./make_deb.py <version> [<path-to-config>]
import sys
sys.path.append('..')
import packager

if __name__ == '__main__':
    version = '0.0.1'
    version = sys.argv[1]
    config_path = None
    if len(sys.argv) > 2:
        config_path = sys.argv[2]

    deb_builder = packager.DebPackage('subshard',
                                      maintainer='Twitchyliquid64 <twitchyliquid64@ciphersink.net>',
                                      description='Subshard is an isolated chrome instance that tunnels all its traffic through a proxy.',
                                      desktop_file='subshard.desktop',
                                      desktop_file_path='subshard.desktop',
                                      bin_files={'../../client/subshard.py': 'subshard',
                                                 '../../client/subshard_configurator.py': 'subshard_configurator',
                                                 'subshard-install-cert.sh': 'subshard-install-cert'},
                                      data_files={'../../client/cr_theme': 'cr_theme',
                                                    'chromeball_google_chrome_poke_by_azerik92-d4c31vz.png': 'chromeball_google_chrome_poke_by_azerik92-d4c31vz.png'},
                                      config_data={'theme_dir': '/usr/share/subshard/cr_theme'},
                                      depends=['openssl', 'libnss3-tools'])

    print deb_builder.package(version, config_path)
