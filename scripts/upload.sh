#!/bin/bash -e

# export GITHUB_TOKEN=xxxx

if [ -z "$1" ]; then
	echo "Usage: ${0} VERSION PATH_TO_ARTIFACTS [--replace]" >> /dev/stderr
	exit 255
fi

if [ -z "$2" ]; then
	echo "Usage: ${0} VERSION PATH_TO_ARTIFACTS [--replace]" >> /dev/stderr
	exit 255
fi

if [ -z $3 ]; then

  ghr -u nildev -r watcher $1 $2

else

  ghr -u nildev -r watcher --replace $1 $2

fi