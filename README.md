Bake
====

Sean Kelleher
-------------

### About

A simple project generator.

### Building

Update GOPATH with a path to the Bake directory. Run `make` to install `bake` to
Bake/bin.

    > make

Create a BAKE environment variable which is set to the path of the Bake
directory, and update your PATH to include the Bake/bin directory. Run the
following for usage information, and to ensure that it is running correctly.

    > bake -h

### Usage

To generate a project with `bake` must specify at least the name of the project,
the name of the owner of the project, and the programming language the owner is
to be written in. An example command that could be used to create the Bake
project itself might be 

    > bake -n Bake -l go -o 'Sean Kelleher'

Where -n specifies the name of the project, -l specifies the language, and -o
specifies the name of the owner of the project. Such an invocation results in a
fairly boring output though:

    Bake/
    Bake/README.md

Only the project directory and a default README have been created, as there is
no indication to `bake` as to what type of project is to be generated. For
instance, running bake with the same parameters but with the addition of a -t
option with a value of bin, results in a template for a Go project which
produces a command-line executable.

    > bake -n Bake -l go -o 'Sean Kelleher' -t bin
    Bake/src/
    Bake/src/bake/
    Bake/src/bake/bake.go

The source directories have now been added with default code so that it may be
compiled and run immediately by using the language's build conventions, in this
case by adding the Bake directory to GOPATH and running

    > go install bake

This will build and output `bake` to a new `bin` directory in Bake, which can be
run with

    > bin/bake
    Bake (C) 2013 Sean Kelleher

### Configuration

Bake is designed so that it doesn't keep any persistent state, so that means
there's no option for configuration of the tool, like storing your name, etc.
Emulating such behaviour is trivial by creating a script called `bake` located
on your path with content such as the following

    #!/bin/bash

    $BAKE/bin/bake -o "Sean Kelleher" -e ezanmoto@gmail.com "$@"

Such a script will allow you to "store" your settings in a script of your own
design, and means that you won't have to go searching for where configuration
files are stored and adhere to configuration formatting, especially since Bake
isn't a tool that one will be using very often. Using scripts like this to
invoke `bake` can be especially handy if you generally only code in one
language, and generally create the same type of project, such as a Java project
that uses Ant and JUnit:

    #!/bin/bash

    $BAKE/bin/bake \
        -o "Sean Kelleher" \
        -e ezanmoto@gmail.com \
        -l java \
        -t ant \
        -t junit \
        "$@"
