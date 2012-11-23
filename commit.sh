#!/bin/sh

# This script should be run to commit items to git repository.

# This script should not be duplicated for running on specific platforms or
# written in a more robust language like perl or python. Shell script has enough
# features to achieve the desired operations and is supported or can be emulated
# on multiple platforms, which makes it as good as an interpreted programming
# language for the purpose of this script.

# Extra build checks
go tool vet src
if [ $? -ne 0 ]
then
    exit 1
fi

# Build
go install bake
if [ $? -ne 0 ]
then
    exit 1
fi

# Run tests
for PROJ in bake tests/perm
do
	go test $PROJ
	if [ $? -ne 0 ]
	then
	    exit 1
	fi
done

# Check formatting
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

# Commit
if [ $# -ne 1 ]
then
	echo "\n$0 expects commit message"
	exit 1
fi

git commit -m "$1"
