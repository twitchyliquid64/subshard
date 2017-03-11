#!/bin/bash
#
# description: Subshard server
#

# Start the service subshard-serv
start() {
        /usr/bin/subshard-serv &
        echo
}

# Restart the service FOO
stop() {
        killall -2 subshard-serv
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
