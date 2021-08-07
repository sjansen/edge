package inspector

import (
	"go/ast"
)

func isContext(expr ast.Expr) bool {
	if sel, ok := expr.(*ast.SelectorExpr); ok {
		if sel.Sel.Name != "Context" {
			return false
		}
		if ident, ok := sel.X.(*ast.Ident); ok {
			return ident.Name == "context"
		}
	}
	return false
}

func isError(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "error"
	}
	return false
}

func matchStruct(expr ast.Expr) (string, string, bool) {
	if star, ok := expr.(*ast.StarExpr); ok {
		if ident, ok := star.X.(*ast.Ident); ok {
			if decl, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
				if _, ok := decl.Type.(*ast.StructType); ok {
					return "", ident.Name, true
				}
			}
		}
		if sel, ok := star.X.(*ast.SelectorExpr); ok {
			return sel.X.(*ast.Ident).Name, sel.Sel.Name, true
		}
	}
	return "", "", false
}
