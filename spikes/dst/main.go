package main

import (
	"fmt"
	"go/token"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

var code = `package main

// magic: enabled
type Foo struct { }

func (f *Foo) Bar() bool {
	return true
}

type (
    // magic: included
	Bar struct { }
	Baz struct { }
)

func main() {
	f := &Foo{}
	b := f.Bar()
}`

func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

type visitor struct{}

func (v *visitor) Visit(node dst.Node) dst.Visitor {
	switch node.(type) {
	case *dst.File:
		return &fileVisitor{}
	case *dst.GenDecl:
		return &declVisitor{}
	}
	return nil
}

type declVisitor struct {
	decs dst.Decorations
}

func (v *declVisitor) Visit(node dst.Node) dst.Visitor {
	switch x := node.(type) {
	case *dst.TypeSpec:
		v.inspectTypeSpec(x)
	}
	return nil
}

type fileVisitor struct{}

func (v *fileVisitor) Visit(node dst.Node) dst.Visitor {
	switch x := node.(type) {
	case *dst.GenDecl:
		if x.Tok == token.TYPE {
			if x.Lparen {
				return &declVisitor{}
			}
			return &declVisitor{
				decs: x.Decorations().Start,
			}
		}
	}
	return nil
}

func (v *declVisitor) inspectTypeSpec(ts *dst.TypeSpec) {
	if len(v.decs) > 0 {
		fmt.Println(v.decs)
	} else {
		fmt.Println(ts.Decorations().Start)
	}
	fmt.Println(ts.Name)
}

func main() {
	fset := token.NewFileSet()
	f, err := decorator.ParseFile(fset, "example.go", code, 0)
	if err != nil {
		die(err)
	}

	dst.Walk(&visitor{}, f)
}
