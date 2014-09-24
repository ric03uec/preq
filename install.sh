#!/bin/bash -e

sys_has() {
    type "$1" > /dev/null 2>&1
    return $?
}
go_download() {
  if sys_has "curl"; then
    curl $*
  elif sys_has "wget"; then
    ARGS=$(echo "$*" | sed -e 's/--progress-bar /--progress=bar /' \
                           -e 's/-L //' \
                           -e 's/-I /--server-response /' \
                           -e 's/-s /-q /' \
                           -e 's/-o /-O /' \
                           -e 's/-C - /-c /')
    wget $ARGS
  fi
}

mkdir -p $HOME/.bin
go_download
mv tpr $HOME/.bin
echo PATH=$PATH:$HOME/.bin >> $HOME/.bashrc
. $HOME/.bashrc 2>&1 > /dev/null
