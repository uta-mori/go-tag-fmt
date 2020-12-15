package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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

		tags := make([]string, len(s.Fields.List))
		for i, f := range s.Fields.List {
			if f.Tag == nil {
				continue
			}
			tags[i] = unquote(f.Tag.Value)
		}

		for i, tag := range Align(tags) {
			if s.Fields.List[i].Tag == nil {
				continue
			}
			s.Fields.List[i].Tag.Value = quote(tag)
		}
		return true
	})

	format.Node(os.Stdout,fset,file)

}

type structTags struct {
	tags []reflect.StructTag
	length map[string]int
	order []string
}

func newStructTags(tags []string) *structTags {
	st := structTags{
		tags: make([]reflect.StructTag, len(tags)),
		length: map[string]int{},
	}
	for i, tag := range tags{
		if tag == "" {
			continue
		}
		rst := reflect.StructTag(tag)
		for _, match := range reTagName.FindAllStringSubmatch(tag, -1){
			tagname := match[1]
			length, ok := st.length[tagname]
			if !ok {
				st.order = append(st.order, tagname)
			}
			if l := len(tagstr(tagname, rst.Get(tagname))); l > length {
				st.length[tagname] = l
			}
		}
		st.tags[i] = rst
	}
	return &st
}

// aligned - alignedはフィールド番号受け取り、整形済みのタグを返します
func (st *structTags) aligned(index int) string {
	b := new(bytes.Buffer)
	for _, tagname := range st.order {
		var t string
		if value, ok := st.tags[index].Lookup(tagname); ok {
			t = tagstr(tagname, value)
		}
		b.WriteString(t)
		b.WriteString(strings.Repeat(" ", st.length[tagname]-len(t)+1))
	}
	return strings.TrimRight(b.String(), " ")
}

func tagstr(tagname, value string) string {
	return tagname + `:"` + value + `"`
}

// unquote - unquoteはタグから`を削除して返します
func unquote(tag string) string {
	s, err := strconv.Unquote(tag)
	if err != nil {
		panic(err)
	}
	return s
}

// quote - quoteはタグに`をつけて返します
func quote(tag string) string {
	return "`" + tag + "`"
}

// reTagName - reTagNameはタグ名を取得する正規表現です
var reTagName = regexp.MustCompile(`(\w+):`)

func Align(tags []string) []string {
	result := make([]string, len(tags))

	st := newStructTags(tags)
	for i := range tags {
		if tags[i] != "" {
			result[i] = st.aligned(i)
		}
	}
	return result
}
