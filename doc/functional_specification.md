Functional Specification
========================

Sean Kelleher
-------------

### Running the program

    bake -o 'Sean Kelleher' -n Bake -l haskell

Will generate a new Haskell project called 'Bake' whose owner is 'Sean
Kelleher'.

#### Arguments

##### Required

+ --owner       -o  The name of the person who owns the project.
+ --name        -n  The name of the new project.
+ --language    -l  The language the new project will be written in.

##### Additional Project Information

+ --type        -t  The type of project, e.g. command, gui, library, etc. Type
                    defaults to 'command' if nothing is specified.
+ --email       -e  The email address of the owner.

##### Help

+ --help        -h  Prints usage information for the program and an overview of
                    the basic commands.
+ --languages   -L  Prints the different languages that projects can currently
                    be generated in.

##### Miscellaneous

+ --verbose     -v  Prints extra information as the program progresses.

### Missing Features

Features that were considered but ultimately left out are provided here with
reasons for their omission.

#### Globals

The owner and email values to bake are unlikely to change over the course of use
of bake by a single developer, so it would make sense to abstract these away to
a global value store. The advantage of this would be that the developer doesn't
have to include these details every time he runs the tool.

However, bake is not the kind of tool that is going to be executed a lot by a
single developer either. There is very little saved by having these values
stored by the program, but if the user wants to change the values in case of
error, or if using bake on a new machine, he will have to search through
documentation to figure out how. In this case, it is actually easier to just
have the developer enter the values manually each time the tool is run. If
sufficient use is made of the tool, it can always be abstracted away within a
one-line wrapper script, such as the following:

    bake --owner 'Sean Kelleher' --email ezanmoto@gmail.com -v "$@"

### Standard Output

The standard output of bake will be the paths of files and directories created
by the tool, separated with newlines. An example might be:

    /home/sean/code/new/doc
    /home/sean/code/new/makefile
    /home/sean/code/new/README
    /home/sean/code/new/src
    /home/sean/code/new/src/main.c
    /home/sean/code/new/src/options.c

This is so that the newly created files may be easily identified and processed.
More information on the progress of the program may be turned on with the
verbosity switch, which will also output files and directories whose generation
have been skipped. An example run of the previous project generation might look
like the following:

    Directory '/home/sean/code/new' already exists, skipping...
    Creating directory '/home/sean/code/new/doc'...
    Creating file '/home/sean/code/new/makefile'...
    Creating file '/home/sean/code/new/README'...
    File '/home/sean/code/new/TODO' already exists, skipping...
    Creating directory '/home/sean/code/new/src'...
    Creating file '/home/sean/code/new/src/main.c'...
    Creating file '/home/sean/code/new/src/options.c'...
