package inspector

import (
	"context"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/types/typeutil"
)

const prefix = "//edge:"

func Inspect(ctx context.Context, basedir string, patterns ...string) ([]*Endpoint, error) {
	if len(patterns) < 1 {
		patterns = []string{"./..."}
	}
	pkgs, err := packages.Load(
		&packages.Config{
			Context: ctx,
			Dir:     basedir,
			Mode:    packages.NeedDeps | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		},
		patterns...,
	)
	if err != nil {
		return nil, err
	}

	i := &inspector{
		endpoints: []*Endpoint{},
	}
	for _, pkg := range pkgs {
		if pkg.IllTyped {
			return nil, ErrTypeErrorsFound
		}
		i.addEndpoints(pkg)
	}

	return i.endpoints, nil
}

type inspector struct {
	endpoints []*Endpoint
}

func (i *inspector) addEndpoint(obj types.Object, routes []string) error {
	endpoint := &Endpoint{
		Package: obj.Pkg().Path(),
		Handler: obj.Name(),
		Routes:  routes,
	}
	i.endpoints = append(i.endpoints, endpoint)
	for _, m := range typeutil.IntuitiveMethodSet(obj.Type(), nil) {
		switch name := m.Obj().Name(); name {
		case "Delete", "Get", "Head", "Options", "Patch", "Post", "Put":
			signature := m.Type().(*types.Signature)
			endpoint.addMethod(name, signature)
		}
	}
	return nil
}

func (i *inspector) addEndpoints(pkg *packages.Package) error {
	for _, file := range pkg.Syntax {
		for _, x := range file.Decls {
			if gd, ok := x.(*ast.GenDecl); ok {
				if gd.Tok == token.TYPE {
					for _, spec := range gd.Specs {
						if gd.Doc != nil {
							ts := spec.(*ast.TypeSpec)
							routes := i.routes(ts, gd.Doc)
							if len(routes) > 0 {
								obj := pkg.TypesInfo.Defs[ts.Name]
								if err := i.addEndpoint(obj, routes); err != nil {
									return err
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (i *inspector) routes(ts *ast.TypeSpec, cg *ast.CommentGroup) (routes []string) {
	if _, ok := ts.Type.(*ast.StructType); ok {
		for _, comment := range cg.List {
			t := comment.Text
			if strings.HasPrefix(t, prefix) {
				t = t[len(prefix):]
				if strings.HasPrefix(t, "route") {
					routes = append(routes, strings.Fields(t[5:])[0])
				}
			}
		}
	}
	return
}
