#!/bin/bash
#
# description: Subshard server
#

# Get function from functions library
. /etc/init.d/functions

# Start the service subshard-serv
start() {
        /usr/bin/subshard-serv &
        ### Create the lock file ###
        touch /var/lock/subsys/subshard-serv
        echo
}

# Restart the service FOO
stop() {
        killall -2 subshard-serv
        ### Now, delete the lock file ###
        rm -f /var/lock/subsys/subshard-serv
        echo
}

### main logic ###
case "$1" in
  start)
        start
        ;;
  stop)
        stop
        ;;
  status)
        status Subshard
        ;;
  restart|reload|condrestart)
        stop
        start
        ;;
  *)
        echo $"Usage: $0 {start|stop|restart|reload|status}"
        exit 1
esac

exit 0
