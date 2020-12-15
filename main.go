package main

import (
	"fmt"
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
	ast.Inspect(file, func(node ast.Node) bool {
		s, ok := node.(*ast.StructType)
		if !ok {
			return true
		}

		for _, f := range s.Fields.List {
			if f.Tag == nil {
				continue
			}
			fmt.Println(f.Tag.Value)
		}
		return true
	})

}
