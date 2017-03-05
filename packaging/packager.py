import errno
import os
import sys
import json
import shutil
import pprint
import stat
from subprocess import call


CRED      = '\033[91m'
CITALIC   = '\33[3m'
CYELLOW   = '\33[33m'
CBLUE2    = '\33[94m'
CEND      = '\033[0m'

def all_read_execute_bits():
    return stat.S_IRUSR | stat.S_IXUSR |\
           stat.S_IRGRP | stat.S_IXGRP |\
           stat.S_IROTH | stat.S_IXOTH

def all_read_bits():
    return stat.S_IRUSR | stat.S_IXUSR |\
           stat.S_IRGRP | stat.S_IXGRP |\
           stat.S_IROTH

def mkdir_p(path):
    try:
        os.makedirs(path)
    except OSError as exc:  # Python >2.5
        if exc.errno == errno.EEXIST and os.path.isdir(path):
            pass
        else:
            raise


class Package(object):
    """Common data structure and methods for packing distribution files
       for all platforms. It is expected that this class will be used by
       platform specific subclasses."""

    def __init__(self, name, temp_dir, **kwargs):
        self.name = name
        self.temp_dir = temp_dir
        self.maintainer = kwargs.get('maintainer', 'Twitchyliquid64 <twitchyliquid64@ciphersink.net>')
        self.description = kwargs.get('description', '')

        self.configuration_dir = kwargs.get('configuration_dir', 'etc/')
        self.binary_dir = kwargs.get('binary_dir', 'usr/bin')
        self.data_dir = kwargs.get('data_dir', 'usr/share')

        self.config_file_name = kwargs.get('config_file_name', '')
        self.bin_files = kwargs.get('bin_files', {})
        self.data_files = kwargs.get('data_files', {})

    def package(self, version, config_path):
        raise NotImplementedError

    def load_config(self, path):
        c = json.load(open(path))
        if 'configuration_dir' in c:
            self.configuration_dir = c['configuration_dir']
        if 'binary_dir' in c:
            self.binary_dir = c['binary_dir']
        if 'data_dir' in c:
            self.data_dir = c['data_dir']

    def _setup_working_dir(self):
        if os.path.exists(self.temp_dir):
            self._log("Deleting old working dir", action=True)
            for root, dirs, files in os.walk(self.temp_dir):
              for momo in dirs:
                os.chmod(os.path.join(root, momo), stat.S_IWUSR | stat.S_IXUSR | stat.S_IRUSR)
              for momo in files:
                try:
                    os.chmod(os.path.join(root, momo), stat.S_IWUSR | stat.S_IXUSR | stat.S_IRUSR)
                except Exception, e:
                    self._log('WARN: Could not chmod() %s: %s', (momo, e))
            shutil.rmtree(self.temp_dir)
        self._log("Creating working dir: %s", self.temp_dir, action=True)
        mkdir_p(self.temp_dir)

    def _setup_package_dirs(self):
        self._log("Make config dir: %s -> %s", (self.configuration_dir, os.path.join(self.temp_dir, self.configuration_dir)), action=True)
        mkdir_p(os.path.join(self.temp_dir, self.configuration_dir))

    def _copy_set(self, file_set, base_dest_path, perms):
        self._log("Make dir: %s", base_dest_path, action=True)
        mkdir_p(base_dest_path)
        for local_path in file_set:
            abs_dest_path = os.path.join(base_dest_path, file_set[local_path])
            is_dir = os.path.isdir(local_path)
            self._log("Copy: %s -> %s (%s)", (local_path, abs_dest_path, 'dir' if is_dir else 'file'), action=True)
            if is_dir:
                shutil.copytree(local_path, abs_dest_path)
            else:
                shutil.copyfile(local_path, abs_dest_path)
            os.chmod(abs_dest_path, perms)

    def _copy_data_and_bin_files(self):
        self._copy_set(self.bin_files, os.path.join(self.temp_dir, self.binary_dir), all_read_execute_bits())
        self._copy_set(self.data_files, os.path.join(self.temp_dir, self.data_dir), all_read_bits())

    def _log(self, message, substitutions=(), **kwargs):
        out = ''
        if kwargs.get('action', False):
            message = CYELLOW + message + CEND
            message = message.replace('%s', CEND + CITALIC + '%s' + CEND + CYELLOW)
        print message % substitutions

    def _log_block(self, content):
        print CBLUE2 + '\t' + content.replace('\n', '\n\t') + CEND

    def _log_object(self, obj):
        content = pprint.PrettyPrinter(indent=4).pformat(obj)
        print CBLUE2 + '\t' + content.replace('\n', '\n\t') + CEND


class OSXPackage(Package):
    plist_xml = '''
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleGetInfoString</key>
  <string>{bundle_info_string}</string>
  <key>CFBundleExecutable</key>
  <string>bin/{executable}</string>
  <key>CFBundleIdentifier</key>
  <string>{unique_identifier}</string>
  <key>CFBundleName</key>
  <string>{name}</string>
  <key>CFBundleIconFile</key>
  <string>{icon}</string>
  <key>CFBundleVersion</key>
  <string>{version}</string>
  <key>CFBundleInfoDictionaryVersion</key>
  <string>6.0</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
</dict>
</plist>
'''


    def __init__(self, name, **kwargs):
        super(OSXPackage, self).__init__(name, os.path.join('/tmp', name + '_osx'), **kwargs)
        self.configuration_dir = kwargs.get('configuration_dir', '%s.app/Contents/configuration' % name)
        self.binary_dir = kwargs.get('binary_dir', '%s.app/Contents/MacOS/bin' % name)
        self.data_dir = kwargs.get('data_dir', '%s.app/Contents/Resources' % name)

        self.bundleInfo = kwargs.get('bundle_info', 'twitchyliquid64 - Copyright (C) 2017')
        self.executable = kwargs.get('executable', self.bin_files.values()[0])
        self.unique_identifier = kwargs.get('unique_identifier', 'com.github.twitchyliquid64.' + name)
        self.icon = kwargs.get('icon', name + ".png")

    def package(self, version, config_path):
        self.version = version
        if config_path:
            self.load_config(config_path)

        self._setup_working_dir()
        self._setup_package_dirs()
        self._copy_data_and_bin_files()

        self._make_plist_file()
        self._make_icon(self.icon)
        os.symlink('/Applications', os.path.join(self.temp_dir, 'Applications'))
        return self._build()

    def _construct_plist(self):
        return OSXPackage.plist_xml.format(
            version=self.version,
            name=self.name,
            bundle_info_string=self.bundleInfo,
            executable=self.executable,
            unique_identifier=self.unique_identifier,
            icon=self.name + '.icns',
        ).strip()

    def _make_plist_file(self):
        plist = self._construct_plist()
        plist_path = os.path.join(self.temp_dir, '%s.app'% (self.name), 'Contents')
        mkdir_p(plist_path)
        self._log("\nWriting Info.plist file to %s", os.path.join(plist_path, 'Info.plist'), action=True)
        self._log_block(plist)
        with open(os.path.join(plist_path, 'Info.plist'), 'w') as outfile:
            outfile.write(plist)

    def _make_icon(self, icon):
        call(['png2icns', os.path.join(self.temp_dir, self.data_dir, self.name + '.icns'), self.icon])

    def _build(self):
        dmg_path = os.path.join(self.temp_dir, '%s.dmg'% (self.name))
        call(['genisoimage', '-V', self.name, '-D', '-R', '-apple', '-no-pad', '-o', dmg_path, self.temp_dir])
        return dmg_path


class DebPackage(Package):
    def __init__(self, name, **kwargs):
        super(DebPackage, self).__init__(name, os.path.join('/tmp', name + '_deb'), **kwargs)

        self.configuration_dir = kwargs.get('configuration_dir', os.path.join('etc/', name))
        self.binary_dir = kwargs.get('binary_dir', 'usr/bin')
        self.data_dir = kwargs.get('data_dir', os.path.join('usr/share', name))

        self.config_file_name = kwargs.get('config_file_name', name + '.json')
        self.desktop_file = kwargs.get('desktop_file', None)
        self.desktop_file_path = kwargs.get('desktop_file_path', self.desktop_file)
        self.config_data = kwargs.get('config_data', {})

        self.architecture = kwargs.get('arch', 'all')
        self.postinst = kwargs.get('postinst', None)
        self.depends = kwargs.get('depends', None)

    def package(self, version, config_path):
        self.version = version
        if config_path:
            self.load_config(config_path)

        self._setup_working_dir()
        self._setup_package_dirs()
        self._copy_data_and_bin_files()

        if self.desktop_file:
            self._make_desktop_entry()
        self._make_control_file()
        self._construct_config_file()

        if self.postinst:
            self._copy_postinst()

        return self._build()

    def _copy_postinst(self):
        postinst_out_path = os.path.join(self.temp_dir, 'DEBIAN', 'postinst')
        self._log("Apply postinst: %s -> %s", (self.postinst, postinst_out_path), action=True)
        shutil.copyfile(self.postinst, postinst_out_path)
        os.chmod(postinst_out_path, stat.S_IWUSR | stat.S_IRUSR | stat.S_IRGRP | stat.S_IROTH | stat.S_IXUSR | stat.S_IXGRP | stat.S_IXOTH)

    def _construct_control_file(self):
        out  = 'Package: %s\n' % self.name
        out += 'Version: %s\n' % self.version
        out += 'Architecture: %s\n' % self.architecture
        if self.depends:
            out += 'Depends: %s\n' % ','.join(self.depends)
        out += 'Maintainer: %s\n' % self.maintainer
        out += 'Description: %s\n' % self.description
        return out

    def _make_control_file(self):
        control = self._construct_control_file()
        mkdir_p(os.path.join(self.temp_dir, 'DEBIAN'))
        self._log("\nWriting control file to %s", os.path.join(self.temp_dir, 'DEBIAN/control'), action=True)
        self._log_block(control)
        with open(os.path.join(self.temp_dir, 'DEBIAN/control'), 'w') as outfile:
            outfile.write(control)

    def _construct_desktop_file(self, path):
        out = ''
        with open(path, "r") as fin:
            for line in fin:
                out += line.replace('VERSION_HERE', self.version).replace('EXEC_PATH_HERE', os.path.join('/', self.binary_dir, self.name))
        return out

    def _make_desktop_entry(self):
        desktop_pkg_path = os.path.join('usr/share/applications/', self.desktop_file)
        self._log('Making desktop entry at %s', desktop_pkg_path, action=True)
        mkdir_p(os.path.join(self.temp_dir, 'usr/share/applications'))
        with open(os.path.join(self.temp_dir, desktop_pkg_path), 'w') as outfile:
            outfile.write(self._construct_desktop_file(self.desktop_file_path))
        os.chmod(os.path.join(self.temp_dir, desktop_pkg_path), stat.S_IWUSR | stat.S_IRUSR | stat.S_IRGRP | stat.S_IROTH)

    def _construct_config_file(self):
        self._log("Make config file: %s", self.config_file_name, action=True)
        conf = self.config_data
        self._log_object(conf)
        with open(os.path.join(self.temp_dir, self.configuration_dir, self.config_file_name), 'w') as outfile:
            json.dump(conf, outfile, sort_keys=True, indent=4, separators=(',', ': '))
            self._log("\t-> %s", os.path.join(self.temp_dir, self.configuration_dir, self.config_file_name), action=True)

    def _build(self):
        call(['dpkg-deb', '--build', self.temp_dir])
        return self.temp_dir + '.deb'
