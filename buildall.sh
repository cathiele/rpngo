#!/bin/bash
set -e
set -x
cd bin
base=$(pwd)
cd $base/ncurses/rpn && go build
cd $base/minimal/rpn && go build
cd $base/tinygo/serialonly && tinygo build -target pico


