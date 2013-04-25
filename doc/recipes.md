Recipes
=======

Sean Kelleher
-------------

### What are Recipes?

A recipe is just a language extension for bake. For instance, python-recipe adds
the functionality for bake to generate Python projects.

### Recipe Contents

Recipes reside in the `recipes` directory and have the following structure:

    {Language}/
        templates/
            ...
        types/
            base
            ...
        tests/
            base
            ...

`templates` contains the templates used to generate projects.

`types` contains project type descriptions describing the different project
types that can be generated.

`tests` contains test scripts for each project type.

### Test Scripts

Tests, as usual, are of 3 critical values:

1. ensure quality
2. show how to use the tool
3. show what to expect

#### Name

The name of a test script is an underscore-separated list of types that this
test checks. For instance, the tests in a test script named `bin_make_test` will
be run before and after a call to 

    bake -o <Owner> -l <Language> -n <ProjectName> -t bin,make,test

#### Structure

Each test script is broken into tests of the following format:

    test one description
    ?command
    ...

    test two description
    ?command
    ...

    ...

    test n description
    ?command
    ...

Each test is separated by a single blank line. A test starts with a description,
followed by 1 or more test actions. A test action is a test directive followed
immediately by a command (one that would be entered into a UNIX shell). The test
directive specifies the expected outcome of the following command, e.g. whether
it should run, whether it should succeed or fail, etc.

The file as a whole is considered to be a test group, which consists of the bake
project types specified by the file name and the tests specified within the
file.

The file should end with a newline, i.e. the last line (that which exists
between the final `\n` and EOF) should be empty.

#### Test Directives

Test directives denote what the test expects the outcome of running the command
following it to be.

##### Setup (` `)

The command following this directive is always expected to pass (return 0). It
is used for setting up tests.

##### Pass (`+`)

The command following this directive is expected to fail (return a non-zero
value) before the bake command is run, and pass (return 0) after the bake
command is run.

##### Build Pass (`=`)

The command following this directive is expected to return an error before the
bake command is run, and is expected to pass (return 0) after the bake command
is run.

This directive is primarily for testing the binary output of a project, whose
attempted execution should cause an error before bake is run. This is because
the executable shouldn't exist, which should cause the OS to signal an error
upon attempting to execute it.

##### Fail (`-`)

The command following this directive is not run before the bake command is run,
and is expected to fail (return a non-zero value) after the bake command is run.

##### Debug (`*`)

The output (to any stream) and return value of the command following this
directive is written to standard output. The command isn't considered a test,
and so doesn't fail on returning an unexpected result. It is for debugging
purposes only, and shouldn't be committed to the recipe's repository.

##### Comment (`/`)

The text following this directive is ignored. It is for debugging purposes only,
and shouldn't be committed to the recipe's repository.

#### Execution

Each test group is run before and after running bake with the types specified in
the test group's type list.  A test group is run by executing the commands in
each test in sequence, and comparing the result of running the command to the
test expectation.

If a test group is testing multiple types, the test group for each of those
individual types is run after bake runs. This helps ensure that specifying
multiple project types doesn't change the behaviour of using bake with
individual types.
