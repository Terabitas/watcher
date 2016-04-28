#!/bin/bash -e

if [ -z "$1" ]; then
	export CONFIG="watcher.conf"
else
	export CONFIG=$1
fi

./bin/watcher --config $CONFIG