#!/bin/bash

PROJECT="store_server_http store_server_rpc"
PROJECT_PDIR=/data/apps/

function usage() {
    echo 'usage: /bin/bash ${PROJECT_PDIR}/[store_server_http|store_server_rpc]/deploy/service.sh start|stop|restart'
    exit 1
}

function check() {
    proc_num=`ps -ef|grep $1 |wc -l`
    if [ $proc_num -ge 1 ]; then
        return 1
    else    
        return 0
    fi    
}

function stop() {
    echo 'STOPPING...'
    check
    if [[ "$?" -eq "1" ]]; then
        service $1 restart
    fi
    echo -e "stop [\033[0;32;32mOK\033[m]"
}

function start() {
    echo 'STARTING...'
    check
    if [[ "$?" -eq "1" ]]; then
       service $1 stop 
    else    
        service $1 start
    fi
    sleep 1
    check
    if [[ "$?" -eq "1" ]]; then
        echo -e "start [\033[0;32;32mOK\033[m]"
        exit 0
    else
        echo "warning: check error, retry..."
        start
    fi
    echo -e "start [\033[0;32;32mOK\033[m]"    
}

function restart() {
    echo 'RESTARTING...'
    stop
    start
}

if test $# -eq 0; then
    usage;
    exit 0;
fi;

case $1 in
    start)
        start;
        ;;
    stop)
        stop;
        ;;
    restart)
        restart;
        ;;
    *)
        usage;
        ;;
esac        
