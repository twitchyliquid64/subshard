#!/usr/bin/python
import os
import platform
import json
import sys
from subprocess import Popen, call

#defaults, populated to platform-specific values in init()
path_to_script_dir = os.path.dirname(os.path.realpath(__file__))
config_dir = '/etc/subshard'
user_config_dir = os.path.join(os.environ['HOME'], '.config', 'subshard')
chrome_path = '/opt/google/chrome/chrome'
proxy_addr = 'https://localhost:8080'
theme_dir = os.path.join(path_to_script_dir, 'cr_theme')
chrome_args = ['--no-first-run', '--disable-default-apps', '--no-default-browser-check', 'http://subshard/']
data_dir = os.path.join(os.path.expanduser("~"), '.subshard_dir')
terminal_command = 'x-terminal-emulator'
path_to_configurator = os.path.join(path_to_script_dir, 'subshard_configurator')


def load_config(path):
    global theme_dir, chrome_path, chrome_args, data_dir, proxy_addr, user_config_dir
    global terminal_command, path_to_configurator
    if not os.path.exists(path):
        return False

    c = json.load(open(path))
    if 'theme_dir' in c:
        theme_dir = c['theme_dir']
    if 'chrome_path' in c:
        chrome_path = c['chrome_path']
    if 'chrome_args' in c: # Override any options except the data dir
        chrome_args = c['chrome_args']
    if 'additional_args' in c: # Add any options in addition to the defaults
        chrome_args += c['additional_args']
    if 'data_dir' in c:
        data_dir = c['data_dir']
    if 'proxy_addr' in c:
        proxy_addr = c['proxy_addr']
    if 'user_config_dir' in c:
        user_config_dir = c['user_config_dir']
    if 'terminal_command' in c:
        terminal_command = c['terminal_command']
    if 'path_to_configurator' in c:
        path_to_configurator = c['path_to_configurator']
    return True



def launch():
    if not os.path.exists(user_config_dir):
        print "No user configuration folder: Opening config utility."
        call([terminal_command, '-e', path_to_configurator], preexec_fn=os.setsid)
        sys.exit(0)

    args = [chrome_path] + chrome_args + ['--user-data-dir=' + data_dir]
    args.append('--load-extension=' + theme_dir)
    args.append('--proxy-server=' + proxy_addr)
    print args
    Popen(args, preexec_fn=os.setsid)



def init():
    global config_dir, chrome_path, theme_dir, terminal_command
    arch = platform.system()
    if arch == 'Linux':
        pass #Defaults as above

    if arch == 'Darwin':
        chrome_path = r'/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'
        config_dir = os.path.realpath(os.path.join(path_to_script_dir, '../../configuration'))
        theme_dir = os.path.realpath(os.path.join(path_to_script_dir, '../../Resources/cr_theme'))
        terminal_command = 'xterm'

    if os.name == 'nt':
        config_dir = os.path.join(os.environ['ProgramFiles'], 'subshard')

    # Load configs if they exist - prefer values for keys in user config over system config.
    load_config(os.path.join(config_dir, 'subshard.json'))
    load_config(os.path.join(user_config_dir, 'subshard.json'))



if __name__ == "__main__":
    init()
    launch()
