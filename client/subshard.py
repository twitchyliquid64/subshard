#!/usr/bin/python
import os
import platform
import json
from subprocess import Popen

#defaults, populated to platform-specific values in init()
path_to_script_dir = os.path.dirname(os.path.realpath(__file__))
arch = ''
config_dir = '/etc/subshard'
chrome_path = '/opt/google/chrome/chrome'
proxy_addr = 'https://localhost:8080'
theme_dir = os.path.join(path_to_script_dir, 'cr_theme')
chrome_args = ['--no-first-run', '--disable-default-apps', '--no-default-browser-check', 'http://subshard/']
data_dir = os.path.join(os.path.expanduser("~"), '.subshard_dir')

def load_config(path):
    global theme_dir, chrome_path, chrome_args, data_dir, proxy_addr
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
    return True


def launch():
    args = [chrome_path] + chrome_args + ['--user-data-dir=' + data_dir]
    args.append('--load-extension=' + theme_dir)
    args.append('--proxy-server=' + proxy_addr)
    print args
    Popen(args, preexec_fn=os.setsid)


def init():
    global arch, config_dir
    arch = platform.system()
    if arch == 'Linux':
        pass #Defaults as above
    if os.name == 'nt':
        config_dir = os.path.join(os.environ['ProgramFiles'], 'subshard')
    load_config(os.path.join(config_dir, 'subshard.json'))


if __name__ == "__main__":
    init()
    launch()
