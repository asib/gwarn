// Package main
package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	warnPrefix = "//:warning"
)

var (
	file = kingpin.Arg("file", "Path to file to check").Required().File()
)

func main() {
	kingpin.Parse()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, (*file).Name(), nil, parser.ParseComments)
	if err != nil {
		fmt.Println("gwarn:", err)
		return
	}

	for _, grp := range f.Comments {
		for _, comment := range grp.List {
			if strings.HasPrefix(comment.Text, warnPrefix) {
				fmt.Println("gwarn:", strings.TrimPrefix(comment.Text, warnPrefix))
			}
		}
	}
}
