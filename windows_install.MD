### You need

1. Python 2.7 installed on your system.
2. [subshard.py](https://github.com/twitchyliquid64/subshard/blob/master/client/subshard.py) on your computer.
3. Copy `client/` to somewhere on your computer.
4. Update the global variables in the python script to point to the correct directories. Specifically, point:
  * `chrome_path` to the file path of your chrome .exe
  * `theme_dir` to the `cr_theme` directory, as per where you copied the `client/` folder earlier.
  * `extension_dir` to the `subshard_extension` directory, as per where you copied the `client/` folder earlier.
  * `user_config_dir` to a directory where you want your config to go. Must exist.
  * `data_dir` to an arbitrary directory where you want Chrome to store its shit.
  * `proxy_addr` to the address of your server.
5. Install the certificate for your server into your trusted certificates. I don't know how this works on Windows, but worst case you can do it in Chrome after you launch subshard (and it has a certificate error).
