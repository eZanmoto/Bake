Source Language Choice
======================

Sean Kelleher
-------------

### Language Options

The language to use for bake will most likely be either a procedural or Object
Oriented language, and due to its relatively small scale either would be
suitable.

A functional programming language wouldn't be particularly suited to the bake
project. In general, functional languages are most suited to

+ Processing a single flow of data with minimal side-effects. The bake project
    requires little processing of data but reads and writes a lot of files.
+ Writing programs whose algorithmic correctness is required. In the case of
    bake, the execution of the program is fairly straightforward, but it is the
    generated projects whose correctness is more important. This must be tested
    using side-effect checking unit tests.
+ Writing programs whose behaviour can be proven with properties, as opposed to
    being asserted with unit tests. Only a small part of bake's execution can be
    proven with properties, the majority must be checked using unit tests.

The rest of this section will cover potential languages to write bake in.
