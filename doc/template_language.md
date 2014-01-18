Bake Template Language
======================

Sean Kelleher
-------------

### About

This document describes the template language bake uses for generating files.
The template syntax was inspired by the syntax of C# format strings.

This template language is (naturally) oriented towards being useful for the bake
project, so not all its features may be relevant to general string formatting
(particularly the boolean operators allowed in conditionals).

### Dictionary

A dictionary (a map from strings to strings) is assumed to be provided to the
operation (function, command, etc.) which is expanding the template.

### Directives

Template directives are delimited by braces ('{' and '}') and are replaced using
various rules. Whitespace is not allowed in directive tags (a single level of
braces).

### Escapes

`{` begins template directives, and so must be escaped by prefixing `{` to
output the literal character. `}` must be escaped with `}` for consistency. An
opening brace (`{`) that isn't escaped is considered to begin a template
directive, whereas a closing brace (`}`) that isn't escaped is an error.

#### Examples

Input (with a dictionary of {}):

    {{}}

Output:

    {}


Input (with a dictionary of {}):

    {}}

Output:

Error: `{` wasn't escaped, and so a directive is assumed to have begun, but `}`
is not allowed at the beginning of a directive.


Input (with a dictionary of {}):

    {{}

Output

Error: `}` hasn't been escaped.


Input (with a dictionary of {}):

    {}

Output:

Error: For the same reason as the second example.

### Variable

If an opening brace is followed by a letter or number, then this directive must
be a variable. Letters and numbers are then read until a closing brace is
encountered - any other characters (or EOF) cause an error. The name that is
read is used as a key into the dictionary - if it exists, its value replaces the
directive, otherwise it is an error.

#### Examples

Input (with a dictionary of {"name": "Sean"}):

    {name}

Output:

    Sean


Input (with a dictionary of {"name": "Sean"}):

    {{name}}

Output:

    {name}


Input (with a dictionary of {"name": "Sean"}):

    {age}

Output:

Error: "age" is not a key in the dictionary.


Input (with a dictionary of {"name": "Sean"}):

    { name }

Output:

Error: Directive tags cannot contain whitespace.

### Conditional Section

If an opening brace is followed by a '?', then this directive begins or closes a
conditional section. A conditional section has the following form:

    {?if}insert-0{:elseif-1}insert-1...{:elseif-n}insert-n{:}default-insert{?}

`if` and `elseif-1`...`elseif-n` are labels representing conditionals, where a
conditional is a number of variable names combined by boolean operators. The
first conditional that evaluates to true using the rules in the following
section has its corresponding insert section processed and inserted - otherwise
the `default-insert` is processed and inserted if no other conditionals evaluate
to true and the `{:}` directive is encountered.

The "dangling-else" issue that is common in programming languages is
circumvented in this rule by nature of the fact that "elseif", "else" and
"endif" directive tags bind to the most recently started "if".

#### Examples

Input (with a dictionary of {"favFood":"pasta", "name":"Sean"}):

    {?name}
    My name is {name}.
    {:favFood}
    I love {favFood}!
    {:}
    Now you know everything about me!
    {?}
    Isn't that great?

Output:

    My name is Sean.
    Isn't that great?


Input (with a dictionary of {"favFood":"pasta"}):

    {?name}
    My name is {name}.
    {:favFood}
    I love {favFood}!
    {:}
    Now you know everything about me!
    {?}
    Isn't that great?

Output:

    I love pasta!
    Isn't that great?


Input (with a dictionary of {"favFude":"pasta"}):

    {?name}
    My name is {name}.
    {:favFood}
    I love {favFood}!
    {:}
    Now you know everything about me!
    {?}
    Isn't that great?

Output:

    Now you know everything about me!
    Isn't that great?


Input (with a dictionary of {"name":"Sean", "favFood":"pasta"}):

    {?name}
    My name is {name}.
    Isn't that great?

Output:

Error: Conditional hasn't been closed.


Input (with a dictionary of {"favFood":"pasta", "name":"Sean"}):

    {?name}
    My name is {name}.
    {?favFood}
    I love {favFood}!
    {?}
    Now you know everything about me!
    {?}
    Isn't that great?

Output:

    My name is Sean.
    I love pasta!
    Isn't that great?

#### Conditionals

A conditional is either a variable name, a bang (`!`) followed by conditional,
or two conditionals combined with a binary operator (either `&` or `|`). Bang
represents boolean "not", ampersand represents boolean "and", and pipe
represents boolean "or". "Not" has the highest precedence, followed by "or",
followed by "and". The binary operators are left-associative. Parentheses may be
added to conditionals for disambiguation.

The evaluation of conditionals is straightforward - a variable evaluates to true
if a key exists in the dictionary with the same value, "not" inverts the value
of the conditional following it, "and" evaluates to true if the conditional on
both sides of it evaluate to true and "or" evaluates to true if the conditional
on either side of it evaluates to true.

##### Examples

Input (with a dictionary of {"favFood":"pasta", "name":"Sean"}):

    {?name&favFood}
    My name is {name} and I love {favFood}!
    {?}
    Isn't that great?

Output:

    My name is Sean and I love pasta!
    Isn't that great?

### Newlines

It may have been noted above, a conditional section which begins at the start of
a line and finishes at the end of a line will have the following newline skipped
if it isn't processed. The same holds when a line ends with an "elseif" or
"else" directive.

#### Examples

Input (with a dictionary of {}):

    These sentences can be tricky.
    {?worry}
    Don't worry though.
    {?}
    It'll be fine

Output:

    These sentences can be tricky.
    It'll be fine.


Input (with a dictionary of {"worry": ""}):

    These sentences can be tricky.
    {?worry}
    Don't worry though.
    {?}
    It'll be fine

Output:

    These sentences can be tricky.
    Don't worry though.
    It'll be fine.


Input (with a dictionary of {"worry": ""}):

    These sentences can be tricky.
    {?worry}
    Don't worry though.
    {:}
    Good thing you didn't worry.
    {?}
    It'll be fine

Output:

    These sentences can be tricky.
    Don't worry though.
    It'll be fine.
