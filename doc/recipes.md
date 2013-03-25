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
            stderr/
            stdout/
            [2..]/
            base
            ...

`templates` contains the templates used to generate projects.

`types` contains project type descriptions describing the different project
types that can be generated.

`tests` contains test scripts for each project type. The sub-folders `stdout`
and `stderr`, and all the numbered sub-folders, contain text templates that
reflect expected outputs for those streams to be used in test scripts.

### Test Scripts

Tests, as usual, are of 3 critical values:
    1. ensure quality
    2. show how to use tool
    3. show what to expect

#### Structure

Each test file is broken into groups of the following format:

    type1 type2 ... typeN

        # test one description
        command
        ...

        # test two description
        command
        ...

        ...

        # test n description
        command
        ...

`type1 type2 ... typeN` is a space-separated list of types that this test group
tests. For instance, a test group with a type list `bin make test` will be
testing the results of running a bake command of the form

    bake -o <Owner> -l <Language> -n <ProjectName> -t bin,make,test

The test group is broken down into further sub-tests by headings (lines
beginning with `#`). Headings define what a sub-test is testing.

Each non-empty line in a test script must be indented by four spaces (not tabs).
This is so that commands that are have special meaning in tests can be denoted
by replacing the third space with a special character, called a test directive.

Test groups are followed by three blank lines, unless the test group doesn't
have a body, in which case the test group is followed by a single blank line.

#### Test Directives

Test directives denote what the test expects the outcome of the following
command to be.

##### Pass (`+`)

The command following this directive is expected to fail (return a non-zero
value) before the bake command is run, and pass (return 0) after the bake
command is run.

##### Fail (`-`)

The command following this directive not run before the bake command is run, and
is expected to fail (return a non-zero value) after the bake command is run.

##### Output (`>`)

This directive is followed by a list of outputs

#### Execution

Each test group is run before and after running bake with the types specified in
the test group's type list. A test group is run by executing its listed commands
in sequence, which are in turn grouped under test descriptions, which describes
and groups a sequence of commands.

When executing the commands, commands without directives are considered set-up
commands and must return 0 before and after running bake, otherwise the test is
considered unstable and fails.

After bake executes with multiple types, the basic test for each type (i.e. the
test containing that type only) is run before the actual test for that
combination of types. In this way, if 
