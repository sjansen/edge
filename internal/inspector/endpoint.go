package inspector

import (
	"go/ast"
)

type Endpoint struct {
	Package string
	Handler string
	Routes  []string
	Delete  *Method
	Get     *Method
	Head    *Method
	Options *Method
	Patch   *Method
	Post    *Method
	Put     *Method
}

type Method struct {
	Params *Struct
	Result *Struct
}

type Struct struct {
	Package string
	Name    string
}

var verbs = map[string]struct{}{
	"Delete":  {},
	"Get":     {},
	"Head":    {},
	"Options": {},
	"Patch":   {},
	"Post":    {},
	"Put":     {},
}

func (e *Endpoint) addMethod(name string, ft *ast.FuncType) error {
	if _, ok := verbs[name]; !ok {
		return nil
	}

	method := &Method{}
	switch n := ft.Params.NumFields(); n {
	default:
		return ErrInvalidSignature
	case 0:
	case 1:
		if !isContext(ft.Params.List[0].Type) {
			return ErrInvalidSignature
		}
	case 2:
		if !isContext(ft.Params.List[0].Type) {
			return ErrInvalidSignature
		}
		if pkg, name, ok := matchStruct(ft.Params.List[1].Type); ok {
			method.Params = &Struct{
				Package: pkg,
				Name:    name,
			}
		} else {
			return ErrInvalidSignature
		}
	}

	switch n := ft.Results.NumFields(); n {
	default:
		return ErrInvalidSignature
	case 0:
	case 1:
		if !isError(ft.Results.List[0].Type) {
			return ErrInvalidSignature
		}
	case 2:
		if !isError(ft.Results.List[1].Type) {
			return ErrInvalidSignature
		}
		if pkg, name, ok := matchStruct(ft.Results.List[0].Type); ok {
			method.Result = &Struct{
				Package: pkg,
				Name:    name,
			}
		} else {
			return ErrInvalidSignature
		}
	}

	switch name {
	case "Delete":
		e.Delete = method
	case "Get":
		e.Get = method
	case "Head":
		e.Head = method
	case "Options":
		e.Options = method
	case "Patch":
		e.Patch = method
	case "Post":
		e.Post = method
	case "Put":
		e.Put = method
	}

	return nil
}
