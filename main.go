package main

import (
	"fmt"
	"os"

	"github.com/sjansen/edge/internal/inspector"
)

func main() {
	endpoints, err := inspector.Inspect("internal/inspector/testdata/api")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	for _, endpoint := range endpoints {
		fmt.Printf("%s.%s\n", endpoint.Package, endpoint.Handler)
		for verb, method := range map[string]*inspector.Method{
			"DELETE":  endpoint.Delete,
			"GET":     endpoint.Get,
			"HEAD":    endpoint.Head,
			"OPTIONS": endpoint.Options,
			"PATCH":   endpoint.Patch,
			"POST":    endpoint.Post,
			"PUT":     endpoint.Put,
		} {
			if method != nil {
				fmt.Printf("  %s (%s) %s\n", verb, method.Params.Name, method.Result.Name)
			}
		}
	}
}
