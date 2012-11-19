#!/bin/sh

# This script should be run to commit items to git repository.

# This script should not be duplicated for running on specific platforms or
# written in a more robust language like perl or python. Shell script has enough
# features to achieve the desired operations and is supported or can be emulated
# on multiple platforms, which makes it as good as an interpreted programming
# language for the purpose of this script.

# Run build
go install bake
if [ $? -ne 0 ]
then
    exit 1
fi

# Run tests
go test bake
if [ $? -ne 0 ]
then
    exit 1
fi

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
git commit
