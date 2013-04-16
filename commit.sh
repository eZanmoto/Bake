#!/bin/sh

# This script should be run to commit items to git repository.

# This script should not be duplicated for running on specific platforms or
# written in a more robust language like perl or python. Shell script has enough
# features to achieve the desired operations and is supported or can be emulated
# on multiple platforms, which makes it as good as an interpreted programming
# language for the purpose of this script.

# Check arguments
if [ $# -ne 1 ]
then
	echo "\n$0 expects commit message"
	exit 1
fi

if [ ${#1} -gt 50 ]
then
	echo "\nCommit message cannot be more than 50 characters"
	exit 1
fi

# Test build
make build
if [ $? -ne 0 ]
then
	exit 1
fi

# Run tests
for TEST in tests/perm bake bake/proj bake/recipe/test diff readers strio
do
	go test -i $TEST
	go test $TEST
	if [ $? -ne 0 ]
	then
		exit 1
	fi
done

# Commit
git commit -m "$1"
