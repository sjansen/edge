package inspector

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

const prefix = "//edge:"

func Inspect(path string) ([]*Endpoint, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	i := &inspector{}
	for name, pkg := range pkgs {
		i.walkPackage(name, pkg)
	}

	endpoints := make([]*Endpoint, 0, len(i.cache.endpoints))
	for _, endpoint := range i.cache.endpoints {
		if len(endpoint.Routes) > 0 {
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints, nil
}

type inspector struct {
	cache struct {
		pkgname   string
		endpoints map[*ast.TypeSpec]*Endpoint
	}
}

func (i *inspector) getEndpoint(ts *ast.TypeSpec) *Endpoint {
	endpoint := i.cache.endpoints[ts]
	if endpoint != nil {
		endpoint.Package = i.cache.pkgname
		endpoint.Handler = ts.Name.Name
	} else {
		endpoint = &Endpoint{
			Package: i.cache.pkgname,
			Handler: ts.Name.Name,
		}
	}
	return endpoint
}

func (i *inspector) walkFile(name string, file *ast.File) {
	for _, x := range file.Decls {
		switch decl := x.(type) {
		case *ast.FuncDecl:
			i.walkFuncDecl(decl)
		case *ast.GenDecl:
			i.walkGenDecl(decl)
		}
	}
}

func (i *inspector) walkFuncDecl(fd *ast.FuncDecl) {
	switch name := fd.Name.Name; name {
	case "Delete", "Get", "Head", "Options", "Patch", "Post", "Put":
		expr := fd.Recv.List[0].Type
		if star, ok := expr.(*ast.StarExpr); ok {
			expr = star.X
		}
		if ident, ok := expr.(*ast.Ident); ok {
			endpoint := i.getEndpoint(ident.Obj.Decl.(*ast.TypeSpec))
			endpoint.addMethod(name, fd.Type)
		}
	}
}

func (i *inspector) walkGenDecl(gd *ast.GenDecl) {
	if gd.Tok == token.TYPE {
		for _, spec := range gd.Specs {
			if gd.Doc != nil {
				i.walkTypeSpec(spec.(*ast.TypeSpec), gd.Doc)
			}
		}
	}
}

func (i *inspector) walkPackage(name string, pkg *ast.Package) {
	i.cache.pkgname = name
	i.cache.endpoints = map[*ast.TypeSpec]*Endpoint{}
	for name, file := range pkg.Files {
		i.walkFile(name, file)
	}
}

func (i *inspector) walkTypeSpec(ts *ast.TypeSpec, cg *ast.CommentGroup) {
	if _, ok := ts.Type.(*ast.StructType); ok {
		endpoint := i.getEndpoint(ts)
		for _, comment := range cg.List {
			t := comment.Text
			if strings.HasPrefix(t, prefix) {
				t = t[len(prefix):]
				if strings.HasPrefix(t, "route") {
					route := strings.Fields(t[5:])[0]
					endpoint.Routes = append(endpoint.Routes, route)
				}
			}
		}
		i.cache.endpoints[ts] = endpoint
	}
}
