#!/usr/bin/python
# This script is used to administrate settings on subshard server.
import sys, getopt, json, getpass, hashlib
import os, errno
from subprocess import call, Popen, PIPE
import re
import pprint, copy

config_path = '/etc/subshard/subshard-serv.json'


#proc    -> name/id of the process
#id = 1  -> search for pid
#id = 0  -> search for name (default)
def process_exists(proc, id=0):
    ps = Popen("ps -A", shell=True, stdout=PIPE)
    ps_pid = ps.pid
    output = ps.stdout.read()
    ps.stdout.close()
    ps.wait()

    for line in output.split("\n"):
        if line != "" and line != None:
            fields = line.split()
            pid = fields[0]
            pname = fields[3]

            if(id == 0):
                if(pname == proc):
                    return True
            else:
                if(pid == proc):
                    return True
    return False


def help():
    print sys.argv[0] + ' [-c <config-file>] <command> [command-specific arguments]'
    print 'COMMANDS:'
    print '\tadduser <username> - adds a user. Prompts for a password.'
    print '\tdeluser <username> - deletes a user.'
    print '\tclean config - Validates the JSON configuration file, and pretty formats it.'
    print '\tsetpw <username> - Sets the password for a user. Prompts for the password.'
    print '\tsetverbosity <enabled/disabled> - Sets the logging verbosity of the server.'
    print '\tsetlistener <listening-address> - Sets the Listening address of the server. EG: \':8080\''
    print '\tsetauthrequired <true/false> - Sets whether a user/password login is needed to use subshard.'
    print '\tblacklist <rule type> <rule> - Creates a blacklist rule.'
    print '\tunblacklist <rule type> <rule> - Deletes a blacklist rule.'
    print '\tshow config - Prints the configuration of the server.'


def parseArgs():
    global config_path
    try:
        opts, args = getopt.getopt(sys.argv[1:],"h",["config="])
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


_ALLOWED_COMMANDS = ('adduser', 'deluser', 'clean', 'setpw', 'setverbosity', 'set-verbosity', 'setlistener', 'set-listener',
                     'setauthrequired', 'set-auth-required', 'blacklist', 'unblacklist', 'show')

def checkArgs(args):
    if len(args) < 2:
        print 'Err: No command and option specified.'
        help()
        sys.exit(1)

    if args[0] not in _ALLOWED_COMMANDS:
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


def setBooleanOption(option_name, data, args):
    if args[0].upper() in ['TRUE', 'YES', 'ENABLE', 'ENABLED']:
        data[option_name] = True
    elif args[0].upper() in ['FALSE', 'NO', 'DISABLE', 'DISABLED']:
        data[option_name] = False
    else:
        print 'Err: This option can only be enabled/disabled.'
        sys.exit(1)


def sanitizeBlacklistArgs(args):
    if len(args) < 2:
        print 'Err: Need to provide both blacklist entry and type. EG: \'subshard-admin blacklist host example.com\''
        sys.exit(1)
    if args[0] not in ['host', 'host-expression', 'prefix', 'expression']:
        print 'Err: Invalid blacklist type %s. Expected host, host-expression, prefix, or expression' % args[0]
        sys.exit(1)
    args[0] = args[0].replace('expression', 'regexp')
    if args[0].endswith('regexp'):
        try:
            re.compile(args[1])
        except re.error, e:
            print 'Err: Invalid regular expression %s: %s' % (args[1], e)
            sys.exit(1)
    return args


def unblacklist(data, args):
    args = sanitizeBlacklistArgs(args)
    entries = data.get('blacklist', [])
    initial_size = len(entries)
    entries = filter(lambda x: x['type'] != args[0] or x['value'] != args[1], entries)
    data['blacklist'] = entries
    if len(entries) == initial_size:
        print 'Err: Entry does not exist.'
        sys.exit(1)


def blacklist(data, args):
    args = sanitizeBlacklistArgs(args)
    entries = data.get('blacklist', [])
    exists = any([x['type'] == args[0] and x['value'] == args[1] for x in entries])
    if exists:
        print 'Err: Blacklist entry already exists.'
        sys.exit(1)
    entries.append({'type': args[0], 'value': args[1]})
    data['blacklist'] = entries


def show(data, args):
    if args[0] == 'config':
        conf = copy.deepcopy(data)
        for user in conf['users']:
            del user['password']
        print pprint.PrettyPrinter(indent=4).pformat(conf)
        sys.exit(0)
    else:
        print 'Err: Did not recognise component \'%s\' to show.' % args[0]
        sys.exit(1)


def forward(data, args):
    entries = data.get('forwarders', [])
    if len(args) > 2:
        exists = any([x['name'] == args[1] for x in entries])
        print exists



if __name__ == '__main__':
    args = parseArgs()
    checkArgs(args)

    data = json.load(open(config_path, 'r'))
    command, args = args[0], args[1:]

    if command == 'adduser':
        adduser(data, args)
    if command == 'deluser':
        deluser(data, args)
    if command == 'setpw':
        adduser(data, args, True)
    if command == 'setverbosity' or command == 'set-verbosity':
        setBooleanOption('verbose', data, args)
    if command == 'setlistener' or command == 'set-listener':
        data['listener'] = args[0]
    if command == 'setauthrequired' or command == 'set-auth-required':
        setBooleanOption('auth-required', data, args)
    if command == 'blacklist':
        blacklist(data, args)
    if command == 'unblacklist':
        unblacklist(data, args)
    if command == 'show':
        show(data, args)
    if command == 'forward':
        forward(data, args)


    try:
        json.dump(data, open(config_path, 'w'), sort_keys=True, indent=4, separators=(',', ': '))
    except IOError as ioex:
        if ioex.errno != 13: #Permission denied
            raise
        print 'Err: Permission denied. Should you run with sudo?'
        sys.exit(1)

    if process_exists('subshard-serv'):
        call(['killall', '-1', 'subshard-serv'])
