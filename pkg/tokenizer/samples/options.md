# Tests for shelldoc options in fenced code blocks

This one says the exit code of the command does not matter:

```shell {shelldocwhatever}
> false
```

This one specifies that the exit code should be 2:

```shell {shelldocexitcode=2}
> (exit 2)
```

More options may follow.
