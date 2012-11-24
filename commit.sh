#!/bin/sh

# This script should be run to commit items to git repository.

# This script should not be duplicated for running on specific platforms or
# written in a more robust language like perl or python. Shell script has enough
# features to achieve the desired operations and is supported or can be emulated
# on multiple platforms, which makes it as good as an interpreted programming
# language for the purpose of this script.

# Test build
./build.sh
if [ $? -ne 0 ]
then
    exit 1
fi

# Commit
if [ $# -ne 1 ]
then
	echo "\n$0 expects commit message"
	exit 1
fi

git commit -m "$1"
