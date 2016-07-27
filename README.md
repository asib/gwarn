# gwarn
A tool that prints warnings in Go source files.

Built with [kingpin](https://github.com/alecthomas/kingpin).

## Install

To install gwarn, run:

```bash
$ go get -u github.com/asib/gwarn
```

## Usage

To run gwarn on all `.go` files in the current directory:

```bash
$ gwarn
```

To run gwarn on all `.go` files in a specific directory:

```bash
$ gwarn dir /path/to/dir
```

To run gwarn on a specific file:

```bash
$ gwarn file /path/to/file.go
```

If you're still unsure:

```bash
$ gwarn help
```

## Example

```go
// main.go
package main

func main() {
  //:warning This is a warning that gwarn will see
}
```

You can then run:

```bash
$ gwarn file main.go
```

and gwarn will output:

```bash
/path/to/file/main.go:5: This is a warning that gwarn will see
```
