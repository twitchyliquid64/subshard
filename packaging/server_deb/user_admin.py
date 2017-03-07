#!/usr/bin/python
# This script is used to administrate users in subshard.
import sys, getopt, json, getpass, hashlib
import os, errno
from subprocess import call

config_path = '/etc/subshard/subshard-serv.json'

def help():
    print sys.argv[0] + ' [-c <config-file>] <command> [command-specific arguments]'

def parseArgs():
    global config_path
    try:
        opts, args = getopt.getopt(sys.argv,"ch",["config="])
        for opt, arg in opts:
            if opt in ("-c", "--config"):
                config_path = arg
            elif opt in ("-h", "--help"):
                help()
                sys.exit()
        return args
    except getopt.GetoptError:
        help()
        sys.exit(2)


def checkArgs(args):
    if len(args) < 2:
        print 'Err: No command and option specified.'
        help()
        sys.exit(1)

    if args[0] not in ('adduser', 'deluser', 'clean', 'setpw'):
        print 'Err: command not recognized.', args
        help()
        sys.exit(3)


def adduser(data, args, replace_if_exists=False):
    users = data.get('users', [])
    exists = any([x['username'] == args[0] for x in users])
    if exists:
        if replace_if_exists:
            users = filter(lambda x: x['username'] != args[0], users)
        else:
            print 'Err: User already exists.'
            sys.exit(1)
    pwd = hashlib.sha256(getpass.getpass()).hexdigest()
    users.append({'username': args[0], 'password': pwd})
    data['users'] = users


def deluser(data, args):
    users = data.get('users', [])
    size = len(users)
    users = filter(lambda x: x['username'] != args[0], users)
    data['users'] = users
    if len(users) == size:
        print 'Err: User does not exist.'
        sys.exit(1)


if __name__ == '__main__':
    args = parseArgs()[1:]
    checkArgs(args)

    data = json.load(open(config_path, 'r'))
    command, args = args[0], args[1:]

    if command == 'adduser':
        adduser(data, args)
    if command == 'deluser':
        deluser(data, args)
    if command == 'setpw':
        adduser(data, args, True)

    try:
        json.dump(data, open(config_path, 'w'), sort_keys=True, indent=4, separators=(',', ': '))
    except IOError as ioex:
        if ioex.errno != 13: #Permission denied
            raise
        print 'Err: Permission denied. Should you run with sudo?'
        sys.exit(1)

    call(['killall', '-1', 'subshard-serv'])
