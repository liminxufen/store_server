#!/bin/bash

set -ue

function usage(){
    echo 'usage: ./deploy.sh <app> <env> <version>'
}

function exit_msg() {
    echo "$1 not in store_server_http store_server_rpc"
    exit 1
}

if test $# -eq 0; then
    usage
    exit 0
fi


list="store_server_http store_server_rpc"
[[ $list =~ (^|[[:space:]])$1($|[[:space:]]) ]] && echo 'OK' || exit_msg
ansible-playbook deploy.yml -i hosts --extra-vars "project=$1 env=$2 check_out=$3 git_password=$GIT_PASSWORD" -c ssh
