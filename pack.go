package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
)

var extFlag = flag.String("e", "png,jpg", "comma-delimited string of allowed file extensions")
var outFlag = flag.String("o", "out.gob", "output file")
var pathFlag = flag.String("p", "./assets", "path to folder containing files")
var stripFlag = flag.Bool("x", false, "strip extensions in output")
var verboseFlag = flag.Bool("v", false, "enable verbose logging")

// allowedExts is the parsed result of extFlag. We use a map instead of a slice for faster indexing.
var allowedExts map[string]struct{} = make(map[string]struct{})

func Pack(root fs.FS) (map[string][]byte, error) {
	pkg := make(map[string][]byte)
	fs.WalkDir(root, ".", func(fPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking %s: %w", fPath, err)
		}
		if d.IsDir() {
			return nil
		}

		// Check file is whitelisted in allowed exts.
		if _, ok := allowedExts[path.Ext(fPath)]; !ok {
			return nil
		}

		// Read data into package.
		f, err := root.Open(fPath)
		if err != nil {
			return fmt.Errorf("opening %s: %w", fPath, err)
		}
		fb, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("reading %s: %w", fPath, err)
		}

		// Strip extensions if enabled.
		if *stripFlag {
			fPath = strings.TrimSuffix(fPath, path.Ext(fPath))
		}

		pkg[fPath] = fb
		if err := f.Close(); err != nil {
			return fmt.Errorf("closing %s: %w", fPath, err)
		}
		if *verboseFlag {
			fmt.Printf("writing %s...\n", fPath)
		}

		return nil
	})

	return pkg, nil
}

func writePkg(pkg map[string][]byte, out io.Writer) error {
	// Encode package into binary format.
	g := gob.NewEncoder(out)
	if err := g.Encode(pkg); err != nil {
		return err
	}
	return nil
}

func main() {
	// Parse flags and allowed file extensions.
	flag.Parse()
	for _, e := range strings.Split(*extFlag, ",") {
		if e[0] != '.' {
			e = "." + e
		}
		allowedExts[e] = struct{}{}
	}

	// Check root path exists.
	if _, err := os.Stat(*pathFlag); err != nil {
		panic(err)
	}

	root := os.DirFS(*pathFlag)
	pkg, err := Pack(root)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(cwdRelative(*outFlag))
	if err != nil {
		panic(err)
	}
	if err := writePkg(pkg, f); err != nil {
		panic(err)
	}

	fmt.Printf("%d assets packaged.\n", len(pkg))
}

func cwdRelative(rel string) string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path.Join(wd, rel)
}
