// Package main
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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
)

func printWarningsInFile(f *ast.File, fset *token.FileSet) {
	for _, grp := range f.Comments {
		for _, comment := range grp.List {
			if strings.HasPrefix(comment.Text, warnPrefix) {
				pos := fset.Position(comment.Pos())
				fmt.Printf("%s:%d: %s\n",
					pos.Filename, pos.Line,
					strings.TrimSpace(strings.TrimPrefix(comment.Text, warnPrefix)))
			}
		}
	}
}

func parseFile(fpath string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fpath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("gwarn:", err)
		os.Exit(1)
	}

	printWarningsInFile(f, fset)
}

func parseDir(dpath string) {
	fset := token.NewFileSet()
	pkgMap, err := parser.ParseDir(fset, dpath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("gwarn:", err)
		os.Exit(1)
	}

	for _, pkg := range pkgMap {
		for _, f := range pkg.Files {
			printWarningsInFile(f, fset)
		}
	}
}

func main() {
	app.Version(version)
	app.Author(author)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case check.FullCommand():
		if wd, err := os.Getwd(); err == nil {
			parseDir(wd)
		} else {
			fmt.Println("gwarn: error: couldn't get current working directory:", err)
			return
		}
	case file.FullCommand():
		parseFile(*filepath)
	case dir.FullCommand():
		parseDir(*dirpath)
	}
}
