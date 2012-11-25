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

+ --type        -t  The type of project, e.g. command, gui, library, etc.
+ --email       -e  The email address of the owner.
+ --license     -i  The software license the project is covered by.

##### Progress Control

+ --resolve     -r  If files that bake wants to create exist, resolve the
                    differences instead of skipping.

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

#### Names

The only characters allowed in inserts and includes are capital and lowercase
letters, and numbers. Variable and project names in inserts and includes should
be written in pascal case (each term in the name should be lowercase, starting
with an uppercase letter, such as ApplesAndOranges). Names are case-insensitive.

#### Outputting Reserved Characters

`{` characters must be escaped by using `{{`, so that they won't be interpreted
as the beginning of an insert or include. For consistency, `}` must also be
escaped, using `}}`.

The delimiters that were considered for containing inserts and includes were

    { }   [ ]   ( )   | |   < >

The chosen delimiters were selected as their position in code is usually easier
to predict and change than that of the others. They also tend to occur less
frequently in non-C-style languages.

#### Project Variable Inserts

This type of variable inserts the associated required variable, or outputs an
error if the named variable does not exist. The form of a project variable
insert is as follows:

    This is the {ProjectName} project.

##### List of Supported Variable Inserts

This is a list of project variables that you can use in templates that should
always be present. Variables are case-sensitive.

+ Owner
+ LowercaseProjectName
+ ProjectName

#### Variable-Dependent Includes

These types of additions are only included if the referenced variable has been
supplied to the program, such as a maintainer email. Such an include has one of
the following two forms:

    {?MaintainerEmail:email:         {MaintainerEmail}}
    {?MaintainerEmail:descrption:    Email {MaintainerEmail} with any issues.}

    {?MaintainerEmail:
    email:          {MaintainerEmail}
    description:    Email {MaintainerEmail} with any issues.
    }

The single-line form replaces the entire declaration (from just before the
opening brace to just after the closing brace). The multi-line form omits the
first newline that follows the colon and the first newline that follows the
closing brace from the output, if said newline characters exist.

A variable-dependent include that depends on multiple variables may join all
required variables together using `&` as a prerequisite for its inclusion. An
example of this may be:

    {?MaintainerName&MaintainerEmail:
    You may reach the maintainer, {MaintainerName}, at {MaintainerEmail}.
    }

A statement such as that in the example will only be included if
`MaintainerName` and `MaintainerEmail` were provided.

You can only nest variable inserts in the body of a variable-dependent include.

##### List of Supported Optional Variable

You may use the following variables in variable dependent includes. Variables
are case-sensitive.

+ Email
+ Licence

#### Project-Dependent Includes

These types of additions are only included if the project being built is of the
referenced type, such as an executable type. Project-dependent has the same
forms and follows the same rules as variable-dependent includes, except it is
denoted with a preceding `!` character, instead of `?`. An example of its use
is:

    {!Executable&Test:
    build/{ProjectName}.o: src/{ProjectName}.c
        gcc src/{ProjectName}.c -o build/{ProjectName}.o
    }

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
