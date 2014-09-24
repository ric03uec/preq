#!/bin/bash -e

BINARY_PATH="https://raw.githubusercontent.com/ric03uec/tpr/v0.1.0/tpr"

sys_has() {
    type "$1" > /dev/null 2>&1
    return $?
}
go_download() {
  if sys_has "wget"; then
    ARGS=$(echo "$*" | sed -e 's/--progress-bar /--progress=bar /' \
                           -e 's/-L //' \
                           -e 's/-I /--server-response /' \
                           -e 's/-s /-q /' \
                           -e 's/-o /-O /' \
                           -e 's/-C - /-c /')
    wget $ARGS
  else
    echo "No wget found, please install wget or download the binary directly from $BINARY_PATH";
  fi
}

mkdir -p $HOME/.bin
go_download "$BINARY_PATH"
mv tpr $HOME/.bin
chmod +x $HOME/.bin/tpr
echo PATH=$PATH:$HOME/.bin >> $HOME/.bashrc
. $HOME/.bashrc 2>&1 > /dev/null
