#!/usr/bin/python

import os
import platform
import json
import sys
import pprint

path_to_script_dir = os.path.dirname(os.path.realpath(__file__))
user_config_dir = os.path.join(os.environ['HOME'], '.config', 'subshard')

CRED      = '\033[91m'
CITALIC   = '\33[3m'
CYELLOW   = '\33[33m'
CBLUE2    = '\33[94m'
CEND      = '\033[0m'

if __name__ == '__main__':
    if len(sys.argv) > 1:
        user_config_dir = sys.argv[1]

    user_config_path = os.path.join(user_config_dir, 'subshard.json')
    print CRED + "    ~======" + CITALIC + "   Subshard Configurator   " + CRED + "======~" + CEND
    condition_str = ' [EXISTS]' if os.path.exists(user_config_path) else ' [new]'
    print CYELLOW + "Config file: " + CEND + CITALIC + user_config_path + CEND + condition_str

    if not os.path.exists(user_config_path):
        if not os.path.exists(user_config_dir):
            os.makedirs(user_config_dir)
        print ''
        print 'Address is formatted: http(s)://<host>:<port>'
        serv = raw_input("What is the address of the subshard server?: ")
        json.dump({'proxy_addr': serv}, open(user_config_path, 'w'))

    else:
        with open(user_config_path, 'r') as fp:
            conf = json.load(fp)
        content = pprint.PrettyPrinter(indent=4).pformat(conf)
        print CBLUE2 + content + CEND
        print ''
        try:
            while True:
                key = raw_input("Type the name of the key you would like to modify: ")
                if key:
                    value = raw_input("Value: ")
                    if value:
                        conf[key] = value
                        with open(user_config_path, 'w') as fp:
                            json.dump(conf, fp)
        except KeyboardInterrupt:
            print ''
