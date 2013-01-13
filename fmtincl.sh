#!/bin/sh

# Copyright 2012 Sean Kelleher. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

# This script should not be duplicated for running on specific platforms or
# written in a more robust language like perl or python. Shell script has enough
# features to achieve the desired operations and is supported or can be emulated
# on multiple platforms, which makes it as good as an interpreted programming
# language for the purpose of this script.

# This script applies formatting rules to bake include files.

for LANG in templates/*
do
	if [ -f $LANG/*.fmt ]
	then
		rm $LANG/*.fmt
	fi

	for INCL in $LANG/*
	do
		if [ $(basename $INCL) != "{ProjectName}" ]
		then
			head -1 $INCL > $INCL.fmt

			DESCR=$(cat $INCL.fmt)
			if [ "$DESCR" = "" -a $(basename $INCL) != "base" ]
			then
				echo "File description for '$INCL' is empty"
				exit 1
			elif [ ${#DESCR} -gt 50 ]
			then
				echo "File description for '$INCL' is over 50 chars (${#DESCR})"
				exit 1
			fi

			bin/fmtincl $INCL >> $INCL.fmt
			if [ $? -eq 0 ]
			then
				diff $INCL $INCL.fmt >/dev/null
				if [ $? -eq 0 ]
				then
					rm $INCL.fmt
				else
					echo "Formatting '$INCL'..."
					mv $INCL.fmt $INCL
				fi
			else
				rm $INCL.fmt
				exit 1
			fi
		fi
	done
done
