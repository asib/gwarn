package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "0.1.0"
	author  = "Jacob Fenton"

	warnPrefix = "//:warning"
)

var (
	app = kingpin.New("gwarn", "A tool that prints warnings in Go source.")

	check = app.Command("check", "Checks the current directory.").Default()

	file     = app.Command("file", "Check a specific file.")
	filepath = file.Arg("filepath", "Path to file to check").Required().ExistingFile()

	dir     = app.Command("dir", "Check a specific directory")
	dirpath = dir.Arg("dirpath", "Path to directory to check").Required().ExistingDir()

	out io.Writer = os.Stdout
)

func printWarningsInFile(f *ast.File, fset *token.FileSet) {
	for _, grp := range f.Comments {
		for _, comment := range grp.List {
			if strings.HasPrefix(comment.Text, warnPrefix) {
				pos := fset.Position(comment.Pos())
				fmt.Fprintf(out, "%s:%d: %s\n",
					pos.Filename, pos.Line,
					strings.TrimSpace(strings.TrimPrefix(comment.Text, warnPrefix)))
			}
		}
	}
}

func parseFile(fpath string) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fpath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	printWarningsInFile(f, fset)
	return nil
}

func parseDir(dpath string) error {
	fset := token.NewFileSet()
	pkgMap, err := parser.ParseDir(fset, dpath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgMap {
		for _, f := range pkg.Files {
			printWarningsInFile(f, fset)
		}
	}

	return nil
}

func main() {
	app.Version(version)
	app.Author(author)

	var err error

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case check.FullCommand():
		if wd, err := os.Getwd(); err != nil {
			fmt.Fprintln(out, "gwarn: couldn't get current working directory:", err)
			return
		} else {
			err = parseDir(wd)
		}
	case file.FullCommand():
		err = parseFile(*filepath)
	case dir.FullCommand():
		err = parseDir(*dirpath)
	}

	if err != nil {
		fmt.Fprintln(out, "gwarn:", err)
	}
}
