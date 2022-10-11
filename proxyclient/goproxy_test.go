package proxyclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetversions(t *testing.T) {

	versions, err := GetVersions(context.Background(), DefaultParam("gopkg.in/src-d/go-billy.v4", nil))
	assert.NoError(t, err)
	assert.Len(t, versions, 11)
}
func TestLatest(t *testing.T) {

	result, err := GetLatest(context.Background(), DefaultParam("gopkg.in/src-d/go-billy.v4", nil))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "v4.3.2", result.Version)
}
func TestGetInfo(t *testing.T) {
	version := "v4.3.2"
	param := DefaultParam("gopkg.in/src-d/go-billy.v4", &version)
	result, err := GetInfo(context.Background(), param)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, version, result.Version)
}
