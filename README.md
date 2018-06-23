# shelldoc: Test Unix shell commands in Markdown documentation

Markdown is widely used for documentation and README.md files that
explain how to use or build some software. Such documentation often
contains shell commands that explain how to build a software or how to
run it. To make sure the documentation is accurate and up-to-date, it
should be automatically tested. *shelldoc* tests Unix shell commands
in Markdown files and reports the results.

## Basic usage

*shelldoc* parses a Markdown input file, detects the code blocks in
it, executes them and compares their output with the content of the
code block. For example, the following code block contains a command,
indicated by either leading a _$_ or a _>_ trigger character, and an
expected response:

    $ echo Hello
    Hello

Lines in code blocks that begin with a _$_ or a _>_ trigger character
are considered commands. Lines inbetween without those lead characters
are considered the expected response. *shelldoc* will execute these
commands and return whether or not the commands succeeded and the
output matches the specificaton:

	% shelldoc README.md
    SHELLDOC: doc-testing "README.md" ...
     CMD (1): echo Hello                                ?  Hello                      :  PASS (match)
     CMD (2): go get github.com/endocode/shelldoc/c...  ?  ...                        :  PASS (match)
    SUCCESS: 2 tests (2 successful, 0 failures, 0 execution errors)

Note that this example is not executed as a test by *shelldoc*, since
it does not start with a trigger character. Trying to do so would
cause an infinite recursion when evaluating the README.md using
*shelldoc*. Try it :-)

## Installation

The usual way to install *shelldoc* is using `go get`:

	$ go get -u github.com/endocode/shelldoc/cmd/shelldoc
	...

Executing documentation may have side effects. For example, running
this `go get` command just installed the latest version of *shelldoc*
in your system.

## Details and syntax

*shelldoc* supports both simple and fenced code blocks. An ellipsis,
as used in the description on how to install *shelldoc* above,
indicates that all output is accepted from this point forward as long
as the command exits with the expected return code (zero, by default).

The `-v (--verbose)` flags enables additional diagnostic output.

A shell is launched that will execute all shell commands in a single
Markdown file. By default, the user's configured shell is used. A
different shell can be specified using the `-s (--shell)` flag:

    % shelldoc --verbose --shell=/bin/sh README.md
	Note: Using user-specified shell /bin/sh.
	...

*shelldoc* uses
the
[Blackfriday Markdown processor](https://github.com/russross/blackfriday) to
parse Markdown files, and the [pflag](https://github.com/spf13/pflag)
package to parse the command line arguments.

## Contributing

*shelldoc*
is
[free and open source software](https://en.wikipedia.org/wiki/Free_and_open-source_software). Everybody
is invited to use, study, modify and redistribute it. To contribute to
*shelldoc*, feel free to fork it and submit pull requests, or to
submit issues in
the
[*shelldoc* issue tracker](https://github.com/endocode/shelldoc/issues). All
contributions are welcome.

To report a bug, the best way is to submit a Markdown file and a
description of how the Markdown file should be interpreted, and how
*shelldoc* interprets it.

## Authors and license

*shelldoc* was developed
by [Mirko Boehm](https://github.com/mirkoboehm). Commercial support,
if necessary, is provided
by [Endocode](https://endocode.com/).

The command line programs of *shelldoc* are located in the `cmd/`
subdirectory and licensed under the terms of the GPL, version 3. The
reusable components are located in the `pkg/` subdirectory and
licensed under the terms of the LGPL version 3.
