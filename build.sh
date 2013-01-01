#!/bin/sh

# This script runs checks, builds and tests.

# This script should not be duplicated for running on specific platforms or
# written in a more robust language like perl or python. Shell script has enough
# features to achieve the desired operations and is supported or can be emulated
# on multiple platforms, which makes it as good as an interpreted programming
# language for the purpose of this script.

# Fails this script if the last command wasn't successful.
exitOnFailure() {
	if [ $? -ne 0 ]
	then
		exit 1
	fi
}

# Clean
go clean
exitOnFailure

# Check formatting
go install fmtincl
exitOnFailure

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

			fmtincl $INCL >> $INCL.fmt
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

# Extra build checks
go tool vet src
exitOnFailure

# Build
go install bake
exitOnFailure

# Run tests
for PROJ in tests/perm bake bake/proj
do
	go test -i $PROJ
	go test $PROJ
	exitOnFailure
done
