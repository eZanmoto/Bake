#!/bin/sh

# Copyright 2012 Sean Kelleher. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

# This script should not be duplicated for running on specific platforms or
# written in a more robust language like perl or python. Shell script has enough
# features to achieve the desired operations and is supported or can be emulated
# on multiple platforms, which makes it as good as an interpreted programming
# language for the purpose of this script.

# This script checks if formatting rules should be applied to Go source code.

FMT=$(gofmt -d -s src)
if [ "$FMT" != "" ]
then
	"$FMT" | tr \n \\n | less

	APPLY=
	while [ "$APPLY" = "" ]
	do
		echo "Apply changes?[y/n/q]"
		read APPLY
		case $APPLY in
			y) gofmt -s -w src ;;
			n) ;;
			q) exit 0 ;;
			*) APPLY=;;
		esac
	done
fi
