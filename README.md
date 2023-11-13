# go-mommy

Mommy's here to support you when using Go~ ❤️

*Go library to read [cargo-mommy](https://github.com/Gankra/cargo-mommy)'s JSON
file and generate similar prompts.*

## Library

See [pkg.go.dev documentation](https://pkg.go.dev/libdb.so/go-mommy).

## CLI

To install the CLI, run:

```sh
go install libdb.so/go-mommy/cmd/mommy@latest # assumes $GOBIN in $PATH
ln -s $(which mommy) $(which daddy) # optional
```

Usage (`mommy --help`):

```
Usage:
  mommy [options] <response-type>
  daddy [options] <response-type>

  <response-type> := <positive> | <negative>
  <positive>      := positive | + | 0
  <negative>      := negative | - | 1

Examples:
  mommy --mommys-little girl positive
  daddy --daddys-little boy --daddys-pronouns his -

Flags:
      --mommys-emotes strings     what emotes mommy will have~
      --mommys-little strings     what to call you~
      --mommys-pronouns strings   what pronouns mommy will use for themself~
      --mommys-roles strings      what role mommy will have~
  -n, --nsfw                      show NSFW flags and include NSFW responses
  -f, --responses-file string     responses.json file (default: bundled responses.json)
  -S, --seed int                  seed for the random number generator (default 1699884773843756301)
  -s, --stylize                   stylize the output using ANSI escape codes
```
