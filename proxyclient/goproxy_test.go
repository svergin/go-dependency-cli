package proxyclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetversions(t *testing.T) {
	modulename := "gopkg.in/src-d/go-billy.v4"

	versions, err := GetVersions(modulename)
	assert.NoError(t, err)
	assert.Len(t, versions, 11)
}
func TestLatest(t *testing.T) {
	modulename := "gopkg.in/src-d/go-billy.v4"

	result, err := GetLatest(modulename)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "v4.3.2", result.Version)
}
func TestGetInfo(t *testing.T) {
	modulename := "gopkg.in/src-d/go-billy.v4"
	v := "v4.3.2"
	result, err := GetInfo(modulename, v)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, v, result.Version)
}
