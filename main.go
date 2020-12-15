package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main(){
	fset := token.NewFileSet()
	filename := "./testdata/a.go"
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal("Error:", err)
	}
	ast.Print(fset, file)

}
