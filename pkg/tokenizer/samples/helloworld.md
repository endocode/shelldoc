# Test: print "Hello World"

Print "Hello" in a ridiculously complicated way:

	$ export HELLOVAR=Hello
    $ echo $HELLOVAR
	Hello

Now print "World", from the root prompt:

    > echo World
    World

Now print a few lines, but only the first one is compared, because of the ellipsis:

    > echo Hello; echo World
    Hello
    ...

The end.
