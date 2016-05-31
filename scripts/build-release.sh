#!/bin/bash -e

VER="$1"
OS=$2
PROJ=$(git rev-parse --show-toplevel)

if [ -z "$1" ]; then
	echo "Usage: ${0} VERSION OS_LIST" >> /dev/stderr
	exit 255
fi

set -u

function package {
	local target=${1}
	local srcdir="${2}/bin"

	local ccdir="${srcdir}/${GOOS}_${GOARCH}"
	if [ -d ${ccdir} ]; then
		srcdir=${ccdir}
	fi
	local ext=""
	if [ ${GOOS} == "windows" ]; then
		ext=".exe"
	fi

	for bin in watcherd; do
		cp ${srcdir}/${bin} ${target}/${bin}${ext}
		cp ${srcdir}/../watcherd.conf.sample ${target}/watcherd.conf
	done
}


function main {
	cd $PROJ
	mkdir -p $PROJ/release

	for os in ${OS}; do
		export GOOS=${os}
		export GOARCH="amd64"

		./build

		TARGET="watcher-${VER}-${GOOS}-${GOARCH}"
		mkdir ${TARGET}
		package ${TARGET} ${PROJ}

		tar cfz release/${TARGET}.tar.gz ${TARGET}
		rm -rf ${TARGET}

		echo "Wrote release/${TARGET}.tar.gz"
	done
}

main