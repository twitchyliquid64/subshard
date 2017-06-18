#!/usr/bin/python
#This script should be run from inside the packaging/server_deb folder.
#./make_deb.py <version> <arch> [<path-to-config>]
import sys
sys.path.append('..')
import packager

if __name__ == '__main__':
    version = '0.0.1'
    version = sys.argv[1]
    build_architecture = sys.argv[2] #eg: amd64, i386, armhf
    config_path = None
    if len(sys.argv) > 3:
        config_path = sys.argv[3]

    deb_builder = packager.DebPackage('subshard-serv',
                                      maintainer='Twitchyliquid64 <twitchyliquid64@ciphersink.net>',
                                      description='Subshard serv is the serverside equivalent to subshard.',
                                      bin_files={'../../serv/serv': 'subshard-serv',
                                                 'admin.py': 'subshard-admin',
                                                 '../../certgen/certgen': 'subshard-gen-keys'},
                                      data_files={'../../serv/web': 'web'},
                                      #binary_dir='bin', -- usr/bin is good
                                      configuration_dir='etc/subshard',
                                      # data_dir='var/lib/subshard', -- actually I think the default works
                                      sysv_script='sysv_script.sh',
                                      config_data={
                                                       'version': version,
                                                       'resources-location': '/usr/share/subshard-serv',
                                                       'listener': ":8080",
                                                       'TLS': {
                                                            'enabled': True,
                                                            'cert-pem-path': '/etc/subshard/cert.pem',
                                                            'key-pem-path': '/etc/subshard/key.pem',
                                                        },
                                                        'users': [
                                                            {
                                                                'username': 'DEFAULT_USERNAME',
                                                                'password': 'DEFAULT_PASSWORD',
                                                            }
                                                        ],
                                                        'auth-required': True,
                                                   },
                                      arch=build_architecture,
                                      postinst='postinst',
                                      depends=['openssl'])

    print deb_builder.package(version, config_path)
