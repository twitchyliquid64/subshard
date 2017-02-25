#!/usr/bin/python
#This script should be run from inside the packaging/deb folder.
#./make_deb.py <path-to-config> <version>
import errno
import os
import sys
import json
import shutil
import pprint
import stat
from subprocess import call

package_name = 'subshard'
maintainer = 'Twitchyliquid64 <twitchyliquid64@ciphersink.net>'
description = 'Subshard is an isolated chrome instance that tunnels all its traffic through a proxy.'
temp_build_dir = '/tmp/subshard_deb'

configuration_dir = 'etc/subshard'
config_file_name = 'subshard.json'

binary_dir = 'usr/bin/'
bin_files = {'../../client/subshard.py': 'subshard'}

data_dir = 'usr/share/subshard'
data_files = {'../../client/cr_theme': 'cr_theme',
              'chromeball_google_chrome_poke_by_azerik92-d4c31vz.png': 'chromeball_google_chrome_poke_by_azerik92-d4c31vz.png'}

def load_config(path):
    global configuration_dir, binary_dir, data_dir

    c = json.load(open(path))
    if 'configuration_dir' in c:
        configuration_dir = c['configuration_dir']
    if 'binary_dir' in c:
        binary_dir = c['binary_dir']
    if 'data_dir' in c:
        data_dir = c['data_dir']


def make_config_structure():
    conf = {}
    conf['theme_dir'] = os.path.join('/', data_dir, 'cr_theme')
    return conf


def mkdir_p(path):
    try:
        os.makedirs(path)
    except OSError as exc:  # Python >2.5
        if exc.errno == errno.EEXIST and os.path.isdir(path):
            pass
        else:
            raise

def copy_set(file_set, base_dest_path, perms):
    print "Make dir: %s" % base_dest_path
    mkdir_p(base_dest_path)
    for local_path in file_set:
        abs_dest_path = os.path.join(base_dest_path, file_set[local_path])
        is_dir = os.path.isdir(local_path)
        print "Copy%s: %s -> %s" % (' dir' if is_dir else '', local_path, abs_dest_path)
        if is_dir:
            shutil.copytree(local_path, abs_dest_path)
        else:
            shutil.copyfile(local_path, abs_dest_path)
        os.chmod(abs_dest_path, perms)

def all_read_execute_bits():
    return stat.S_IRUSR | stat.S_IXUSR |\
           stat.S_IRGRP | stat.S_IXGRP |\
           stat.S_IROTH | stat.S_IXOTH

def all_read_bits():
    return stat.S_IRUSR | stat.S_IXUSR |\
           stat.S_IRGRP | stat.S_IXGRP |\
           stat.S_IROTH

def make_control_file(version):
    # TODO: add requires for google-chrome-stable > 54
    out  = 'Package: %s\n' % package_name
    out += 'Version: %s\n' % version
    out += 'Architecture: all\n'
    out += 'Maintainer: %s\n' % maintainer
    out += 'Description: %s\n' % description
    return out

def get_desktop_file(version):
    out = ''
    with open("subshard.desktop", "r") as fin:
        for line in fin:
            out += line.replace('VERSION_HERE', version).replace('EXEC_PATH_HERE', os.path.join('/', binary_dir, 'subshard'))
    return out


if __name__ == '__main__':
    version = '0.0.1'
    if len(sys.argv) > 1:
        load_config(sys.argv[1])
        version = sys.argv[2]

    if os.path.exists(temp_build_dir):
        print "Deleting old working dir"
        for root, dirs, files in os.walk(temp_build_dir):
          for momo in dirs:
            os.chmod(os.path.join(root, momo), stat.S_IWUSR | stat.S_IXUSR | stat.S_IRUSR)
          for momo in files:
            os.chmod(os.path.join(root, momo), stat.S_IWUSR | stat.S_IXUSR | stat.S_IRUSR)
        shutil.rmtree(temp_build_dir)

    print "Make working dir: %s" % temp_build_dir
    mkdir_p(temp_build_dir)

    print "\nMake config dir: %s -> %s" % (configuration_dir, os.path.join(temp_build_dir, configuration_dir))
    mkdir_p(os.path.join(temp_build_dir, configuration_dir))

    print "Make config file:"
    conf = make_config_structure()
    pprint.PrettyPrinter(indent=4).pprint(conf)
    with open(os.path.join(temp_build_dir, configuration_dir, config_file_name), 'w') as outfile:
        json.dump(conf, outfile)
        print "    -> %s" % os.path.join(temp_build_dir, configuration_dir, config_file_name)

    print ''
    copy_set(bin_files, os.path.join(temp_build_dir, binary_dir), all_read_execute_bits())
    copy_set(data_files, os.path.join(temp_build_dir, data_dir), all_read_bits())

    print '\nMaking desktop entry at /usr/share/applications/subshard.desktop'
    mkdir_p(os.path.join(temp_build_dir, 'usr/share/applications'))
    with open(os.path.join(temp_build_dir, 'usr/share/applications/subshard.desktop'), 'w') as outfile:
        outfile.write(get_desktop_file(version))
    os.chmod(os.path.join(temp_build_dir, 'usr/share/applications/subshard.desktop'), stat.S_IWUSR | stat.S_IRUSR | stat.S_IRGRP | stat.S_IROTH)

    control = make_control_file(version)
    mkdir_p(os.path.join(temp_build_dir, 'DEBIAN'))
    print "\nWriting control file to %s" % os.path.join(temp_build_dir, 'DEBIAN/control')
    print '\t' + control.replace('\n', '\n\t')
    with open(os.path.join(temp_build_dir, 'DEBIAN/control'), 'w') as outfile:
        outfile.write(control)

    print "Building package"
    call(['dpkg-deb', '--build', temp_build_dir])
