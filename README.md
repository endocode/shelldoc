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

~~~shell
% shelldoc run README.md
SHELLDOC: doc-testing "go/src/github.com/endocode/shelldoc/README.md" ...
 CMD (1): echo Hello                                ?  Hello                      :  PASS (match)
 CMD (2): go get -u github.com/endocode/shelldo...  ?  ...                        :  PASS (match)
 CMD (3): export GREETING="Hello World"             ?  (no response expected)     :  PASS (execution successful)
 CMD (4): echo $GREETING                            ?  Hello World                :  PASS (match)
SUCCESS: 4 tests (4 successful, 0 failures, 0 execution errors)
~~~

Note that this example is not executed as a test by *shelldoc*, since
it does not start with a trigger character. Trying to do so would
cause an infinite recursion when evaluating the README.md using
*shelldoc*. Try it :-) The percent symbol is commonly used as a shell
prompt next to  _$_ or a _>_. It can be used in documentation as a
prompt indicator without triggering a *shelldoc* test.

## Installation

The usual way to install *shelldoc* is using `go get`:

	$ go get -u github.com/endocode/shelldoc/cmd/shelldoc
	...

Executing documentation may have side effects. For example, running
this `go get` command just installed the latest version of *shelldoc*
in your system. Containers or VMs can be used to isolate such side
effects.

## Details and syntax

All code blocks in the Markdown input are evaluated and executed as
tests. A test succeeds if it returns the expected exit code, and the
output of the command matches the response specified in the code
block.

*shelldoc* supports both simple and fenced code blocks. An ellipsis,
as used in the description on how to install *shelldoc* above,
indicates that all output is accepted from this point forward as long
as the command exits with the expected return code (zero, by default).

The `-v (--verbose)` flags enables additional diagnostic output.

A shell is launched that will execute all shell commands in a single
Markdown file. By default, the user's configured shell is used. A
different shell can be specified using the `-s (--shell)` flag:

    % shelldoc --verbose run --shell=/bin/sh README.md
	Note: Using user-specified shell /bin/sh.
	...

The shell's lifetime is that of the test run of a single Markdown
file. The environment of the shell is available between test
interactions:

	$ export GREETING="Hello World"
	$ echo $GREETING
	Hello World

*shelldoc* uses
the
[Blackfriday Markdown processor](https://github.com/russross/blackfriday) to
parse Markdown files, and the [pflag](https://github.com/spf13/pflag)
package to parse the command line arguments.

## Options

Regular code blocks do not have a way to specify options. The only
thing that can be specified about them are the commands and the
responses. That means the expected return code must always be zero for
the test to succeed.

Sometimes, however, things are more complicated. Some commands are
expected to return a different exit code than zero. Some commands
return exit codes that are unknown up-front. Both options can be
handled by specifying tests in fenced code blocks. Fenced code blocks
may have an info string after the opening characters. This info string
is usually used to specify the language of the following code. After
the language specifier however, other information may
follow. `shelldoc` uses this opportunity to allow the user to specify
options about the test. These options are:

	```shell {shelldocwhatever}
    % echo Hello && false
    Hello
    ```
Try executing this test:

```shell {shelldocwhatever}
> echo Hello && false
Hello
```

The _shelldocwhatever_ options tells `shelldoc` that the exit code of
the following command does not matter. If any expected response is
specified, it will still be evaluated.

    ```shell {shelldocexitcode=2}
    % (exit 2)
    ```

This executes the test, for tests:

```shell {shelldocexitcode=2}
> (exit 2)
```

The _shelldocexitcode_ specifies an exact exit code that is
expected. The test fails if the exit code of the command does not
match the specified one, or if the response does not match the
expected response.

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
by [Mirko Boehm](http://www.creative-destruction.org). Commercial support,
if necessary, is provided
by [Endocode](https://endocode.com/).

The command line programs of *shelldoc* are located in the `cmd/`
subdirectory and licensed under the terms of the GPL, version 3. The
reusable components are located in the `pkg/` subdirectory and
licensed under the terms of the LGPL version 3. Unit test and example
code is licensed under the Apache-2.0 license.
