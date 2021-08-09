package inspector

import (
	"go/types"
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

func (e *Endpoint) addMethod(name string, sig *types.Signature) error {
	method := &Method{}

	params := sig.Params()
	switch n := params.Len(); n {
	default:
		return ErrInvalidSignature
	case 0:
	case 1:
		// TODO: also accept struct
		if !isContext(params.At(0).Type()) {
			return ErrInvalidSignature
		}
	case 2:
		if !isContext(params.At(0).Type()) {
			return ErrInvalidSignature
		}
		if pkg, name, ok := matchStruct(params.At(1).Type()); ok {
			if pkg == e.Package {
				method.Params = &Struct{
					Name: name,
				}
			} else {
				method.Params = &Struct{
					Package: pkg,
					Name:    name,
				}
			}
		} else {
			return ErrInvalidSignature
		}
	}

	results := sig.Results()
	switch n := results.Len(); n {
	default:
		return ErrInvalidSignature
	case 0:
	case 1:
		// TODO: also accept struct
		if !isError(results.At(0).Type()) {
			return ErrInvalidSignature
		}
	case 2:
		if pkg, name, ok := matchStruct(results.At(0).Type()); ok {
			if pkg == e.Package {
				method.Result = &Struct{
					Name: name,
				}
			} else {
				method.Result = &Struct{
					Package: pkg,
					Name:    name,
				}
			}
		} else {
			return ErrInvalidSignature
		}
		if !isError(results.At(1).Type()) {
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
