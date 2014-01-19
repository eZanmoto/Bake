Debugging Tests
===============

Sean Kelleher
-------------

### Using `gdb`

To debug tests with `gdb`, first generate the test binaries to `pkg/tst` with

    make tests

To start the test binaries in `gdb`, change to the `pkg/tst` directory (which
contains a script to initialize `gdb` for Go binaries) and supply the
appropriate package`.test` file to the `gdb` command. For example, to debug the
`bake/proj` test, run

    gdb bake/proj.test
