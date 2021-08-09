package inspector

import (
	"go/types"
)

func isContext(typ types.Type) bool {
	if named, ok := typ.(*types.Named); ok {
		obj := named.Obj()
		if obj.Pkg().Path() == "context" {
			return obj.Name() == "Context"
		}
	}
	return false
}

func isError(typ types.Type) bool {
	if named, ok := typ.(*types.Named); ok {
		obj := named.Obj()
		if obj.Pkg() == nil {
			return obj.Name() == "error"
		}
	}
	return false
}

func matchStruct(typ types.Type) (string, string, bool) {
	if ptr, ok := typ.(*types.Pointer); ok {
		if named, ok := ptr.Elem().(*types.Named); ok {
			if _, ok := named.Underlying().(*types.Struct); ok {
				pkg := named.Obj().Pkg().Path()
				name := named.Obj().Name()
				return pkg, name, true
			}
		}
	}
	return "", "", false
}
