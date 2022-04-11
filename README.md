# gobpack

gobpack is a simple tool for packaging a collection of files into a distributable [gob.](https://go.dev/blog/gob)

## Description

I created gobpack to help bundle static assets for [Gelyn's Edge](https://github.com/gelynsedge) into a separate, distributable binary. WebAssembly has restrictions on file access and bundling the assets alongside the game would negatively impact download speeds. This solution, although not perfect, serves my needs for now.

## Installation

gobpack requires Go 1.18 or higher. Also ensure your `$GOPATH/bin` is in the global path.

To install:

```bash
$ go install github.com/kinesivan/gopack@latest
```

## Usage

```bash
$ gobpack -h
Usage of gobpack:
  -e string
        comma-delimited string of allowed file extensions (default "png,jpg")
  -o string
        output file (default "out.gob")
  -p string
        path to folder containing files (default "./assets")
  -v    enable verbose logging
```

```bash
$ gobpack -p ./assets -o assets.gob
```

## License

This project is licensed under the MIT License - see the LICENSE.md file for details
