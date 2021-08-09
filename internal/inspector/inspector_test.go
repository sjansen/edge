package inspector_test

import (
	"context"
	"sort"
	"testing"

	"github.com/sjansen/edge/internal/inspector"
	"github.com/stretchr/testify/require"
)

func TestInspect(t *testing.T) {
	require := require.New(t)

	actual, err := inspector.Inspect(context.TODO(), "testdata")
	require.NoError(err)

	sort.Slice(actual, func(i, j int) bool {
		if actual[i].Package < actual[j].Package {
			return actual[i].Handler < actual[j].Handler
		}
		return actual[i].Package < actual[j].Package
	})

	expected := []*inspector.Endpoint{
		{
			Package: "example.com/test/api",
			Handler: "HelloHandler",
			Routes:  []string{"/hello"},
			Get: &inspector.Method{
				Result: &inspector.Struct{
					Package: "",
					Name:    "HelloResponse",
				},
			},
			Post: &inspector.Method{
				Params: &inspector.Struct{
					Package: "",
					Name:    "HelloParams",
				},
				Result: &inspector.Struct{
					Package: "",
					Name:    "HelloResponse",
				},
			},
		}, {
			Package: "example.com/test/api",
			Handler: "RegistrationHandler",
			Routes:  []string{"/users/register"},
			Post: &inspector.Method{
				Params: &inspector.Struct{
					Package: "example.com/test/shared",
					Name:    "UserProfile",
				},
			},
		},
	}
	require.Equal(expected, actual)
}
