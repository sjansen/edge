package main

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type User struct {
	FirstName string `validate:"alphaunicode,required"`
	LastName  string `validate:"alphaunicode,required"`
	Email     string `validate:"required,email"`
	Age       uint8  `validate:"gte=0,lte=150"`
}

func TestValidator(t *testing.T) {
	require := require.New(t)
	validate := validator.New()

	u := &User{
		FirstName: "Stuart",
		LastName:  "Jansen",
		Email:     "sjansen@example.com",
		Age:       42,
	}
	err := validate.Struct(u)
	require.NoError(err)

	u = &User{
		FirstName: "Stuart",
		LastName:  "Jansen",
		Email:     "6 × 9",
		Age:       42,
	}
	err = validate.Struct(u)
	require.Error(err)
}
