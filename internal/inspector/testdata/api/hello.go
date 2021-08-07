package api

import (
	"context"
	"fmt"
)

type HelloParams struct {
	Name string `validate:"alphaunicode,required"`
}

type HelloResponse struct {
	Message string
}

//edge:route /hello
type HelloHandler struct{}

func (h *HelloHandler) Get(ctx context.Context) (*HelloResponse, error) {
	return &HelloResponse{Message: "Hello, World!"}, nil
}

func (h *HelloHandler) Post(ctx context.Context, params *HelloParams) (*HelloResponse, error) {
	msg := fmt.Sprintf("Hello, %s!", params.Name)
	return &HelloResponse{Message: msg}, nil
}
