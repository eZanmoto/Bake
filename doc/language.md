Source Language Choice
======================

Sean Kelleher
-------------

### Language Options

The language to use for bake will most likely be either a procedural or Object
Oriented language, and due to its relatively small scale either would be
suitable. A functional programming language wouldn't be particularly suited to
the bake project. In general, functional languages are most suited to

+ Processing a single flow of data with minimal side-effects. The bake project
    requires little processing of data but reads and writes a lot of files.
+ Writing programs whose algorithmic correctness is required. In the case of
    bake, the execution of the program is fairly straightforward, but it is the
    generated projects whose correctness is more important. This must be tested
    using side-effect checking unit tests.
+ Writing programs whose behaviour can be proven with properties, as opposed to
    being asserted with unit tests. Only a small part of bake's execution can be
    proven with properties, the majority must be checked using unit tests.

The rest of this section will cover potential languages to write bake in. The
focus will be on the portability of the language since the project itself isn't
particularly large. An ideal language would also be one in which the code can be
expressed eloquently. Since the program isn't likely to be run often, efficiency
of the language isn't a major concern.

#### Eiffel

A very safe and sturdy looking language, probably more industrial-strength and
correct than this project warrants. Similarly to the argument against functional
properties in the use of this project, the correctness and elegance of
invariants would be wasted on bake.

#### Erlang

Like Eiffel, a very safe and robust language, particularly suited to
concurrency. Such capabilities could be appropriate for bake.

#### Go

Compiles to native code, but active support for numerous platforms, including
Linux, Windows, OS X, FreeBSD, and also ships source code for porting. Provides
a straightforward, imperative programming interface with concurrency and a
simple programming style. Would be a very good fit for the bake project.

#### Lua

Very lean and lightweight language; while suitable for the project, a more
robust and maintainable alternative might be preferred.

#### OCaml

Nice and stable looking language with both object oriented and functional
language support. It has a lot more power than is required for a project of this
size.

#### Self

Nice concept, but the lack of static types and classes will not be beneficial
for this project.

#### Smalltalk

A minimal yet functional language; its lean nature makes it tempting for use in
this project, however, it doesn't seem to have a lot of support on many
platforms.

### Source Language

The chosen source language is Go, due to its lightweight syntax and
implementation, but also its portability - it seems to be actively developed on
a variety of platforms.
