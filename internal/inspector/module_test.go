package inspector_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/modfile"
)

func TestModfile(t *testing.T) {
	require := require.New(t)

	data, err := ioutil.ReadFile("testdata/go.mod")
	require.NoError(err)

	f, err := modfile.Parse("testdata/go.mod", data, nil)
	require.NoError(err)
	require.Equal("example.com/test", f.Module.Mod.Path)
}
