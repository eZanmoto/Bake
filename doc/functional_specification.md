Functional Specification
========================

Sean Kelleher
-------------

### Running the program

    bake -o 'Sean Kelleher' -n Bake -l haskell

Will generate a new Haskell project called 'Bake' whose owner is 'Sean
Kelleher'. Using PascalCase for the project name will allow bake to decompose
the name to generate language-conforming directory/file names, as well as
generating meaningful documentation.

#### Arguments

##### Required

+ --owner       -o  The name of the person who owns the project.
+ --name        -n  The name of the new project.
+ --language    -l  The language the new project will be written in.

##### Additional Project Information

+ --type        -t  The type of project, e.g. command, gui, library, etc. It can
                    be specified multiple times, and each specification can
                    contain multiple types, separated by commas.
+ --email       -e  The email address of the owner.
+ --license     -i  The software license the project is covered by.

##### Progress Control

+ --resolve     -r  If files that bake wants to create exist, resolve the
                    differences instead of skipping. This adds new lines to the
                    existing file in a diff-like format.
+ --include     -R  If files that bake wants to create exist, resolve the
                    differences instead of skipping. This adds new lines to the
                    existing file in what is inferred to be the correct
                    position. This option should only be used after a project
                    has been generated, to include a forgotten project type for
                    instance. If a project has been significantly updated,
                    consider using `--resolve` instead.
+ --merge       -m  Like resolve, but the changes are integrated without being
                    highlighted with the diff-like syntax.

##### Help

+ --help        -h  Prints usage information for the program and an overview of
                    the basic commands.
+ --languages   -L  Prints the different languages that projects can currently
                    be generated in.
+ --types       -T  Prints the different project types that can be added to a
                    project in the specified language.
+ --licenses    -I  Prints the different licenses that can be used.

##### Miscellaneous

+ --verbose     -v  Prints extra information as the program progresses.
+ --default     -d  Use the arguments given to options during this run as the
                    default values of those options.
+ --rm-default  -D  Remove the default values of the specified language.
+ --set-default -s  Set default values (as with `--default`) without running
                    tool.

### Missing Features

Features that were considered but ultimately left out are provided here with
reasons for their omission.

#### Support for Mutually Exclusive Types

This section discusses the need for mutually exclusive project types, such as
two projects whose compilation depends on a specific environment or
compiler/interpreter implementation.

This functionality would clutter up the project templates with
implementation-specific (and therefore, not general, or necessarily
widely-adopted) code.

It would complicate the generation process, by adding more rules and checks at
project generation time.

It would add redundancy to files, as every pair of exclusive projects would need
to list each other, possibly resulting in bugs due to one project being listed
as exclusive of another who doesn't think it is, or else, duplicating checks.

It would complicate testing, as some generated projects would have to be tested
with specific compilers and interpreters.

The benefit of this feature (minor convenience for what appears to be a small
subset of use cases) is vastly overshadowed by its drawbacks, and will not be
included.

### Project Templates

Project templates are located in the $BAKE\templates directory. An environment
variable is used to allow access to the templates directory in a platform
independent manner. This will be referred to as the templates directory for the
remainder of this section.

#### Language Template Structure

Templates for projects in specific languages are stored directly in the
templates directory. The layout of such language templates are described as
follows:

    templates/LanguageName/{ProjectName}/
    templates/LanguageName/Base
    templates/LanguageName/Executable
    templates/LanguageName/Library

Where the `{ProjectName}` directory contains all the actual templates. The
`Base` file contains a listing of paths of all templates that are included
regardless of the actual project type. The `Executable` and `Library` files
contain listings of paths of templates that will be included in addition to the
"base" file list, if these project types are specified. Such project include
files have the following form:

    One line description
    file1
    file2
    file3
    directory1/
        file4
        directory2/
            file5
        file6
    file7
    directory3/
        file8

The one line description must be a maximum of 50 characters or an error will be
signalled. The description is output as part of the result of the T command.
This description is not optional, and is not terminated by a period.

Each item ending in `/` denotes a directory and each file and each item at a
particular level is thought to be contained in the first preceding directory at
the higher level.

The reason for this approach is its minimalist yet concise nature, it is
relatively easy to read and parse. The use of indentation removes the need for
listing the directory path for each file separately.

#### Template File Format

See the template\_language.md document for a specification of the template
language used by bake. When generating files, bake passes a dictionary
consisting of the supplied project types mapping to `""` variables mapping
to their values. Project types begin with lowercase letters and variables begin
with uppercase letters to avoid namespace collisions.

### Project Types

The type of project to be generated is supplied to bake using the `--type` or
`-t` options. Each project type specified causes that project type's files
(specified in the language's template root) to be generated, and
project-dependant includes that specify that project to be included.

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

### Examples

#### Default Values

`--default` and `--set` are used to set the current parameters (except for the
project name) as the default values of those parameters for the current
language, although some global options (such as name and email) are set as
default parameter values for all languages. An possible example of their usage
is as follows:

    > bake -o 'Sean Kelleher' -l c -n bake
    /home/sean/code/bake/README
    /home/sean/code/bake/src
    /home/sean/code/bake/src/bake.c

    > bake -d -o 'Sean Kelleher' -l c -t make,lib -n xlib
    /home/sean/code/xlib/makefile
    /home/sean/code/xlib/README
    /home/sean/code/xlib/src
    /home/sean/code/xlib/src/xlib.c
    /home/sean/code/xlib/include
    /home/sean/code/xlib/include/xlib.h

    > bake -l c -n ylib
    /home/sean/code/ylib/makefile
    /home/sean/code/ylib/README
    /home/sean/code/ylib/src
    /home/sean/code/ylib/src/ylib.c
    /home/sean/code/ylib/include
    /home/sean/code/ylib/include/ylib.h

Because `-d` was passed in the second example run of bake, `Sean Kelleher` is
set as the default value for the `-o` parameter (otherwise the example wouldn't
run) and `make,lib` is set as the default value for the `-t` parameter.
