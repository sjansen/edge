package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sjansen/edge/internal/inspector"
)

func main() {
	endpoints, err := inspector.Inspect(context.Background(), "internal/inspector/testdata")
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
				var params, result string
				if method.Params != nil {
					params = method.Params.Name
				}
				if method.Result != nil {
					result = method.Result.Name
				}
				fmt.Printf("  %s (%s) %s\n", verb, params, result)
			}
		}
	}
}
