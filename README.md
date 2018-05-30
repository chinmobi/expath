# go-expath

[![Go Report Card](https://goreportcard.com/badge/github.com/chinmobi/expath)](https://goreportcard.com/report/github.com/chinmobi/expath)

The expath is a Go package that extends the Go standard library path/filepath's Match and Glob routines, especially supporting any-directories' pattern (** matches zero or more directories in a path).

## Usage

```go
matched, err := expath.Match(`/foo/b*/**/z*.txt`, `/foo/begin/a/b/c/zero.txt`)
matches, atRoot, err := expath.Glob(`/foo/b*/**/z*.txt`, `./`)
```

## Installation

```
$ go get github.com/chinmobi/expath
```

## License

MIT

## Author

Zhaoping Yu (<yuzhaoping1970@gmail.com>)
