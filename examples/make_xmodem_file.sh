#!/bin/bash

# Packages the givn list of files so that they can be transferred
# and saved in a single xmodem send
#
# Example:
#
# On PC:
#
# ./make_xmodem_send *.rpn > out.txt
# screen /dev/ttyUSB0 115200
#
# In screen ctrl-a :, then
# exec !! sx out.txt
#
# On calculator
# rx @

function dump_file {
	path=$1
	echo '{'
	cat $path
	echo '}'
	echo "'$(basename $path)' save"
}

for f in $@; do
	dump_file $f
done

