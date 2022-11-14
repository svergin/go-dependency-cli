package goproxy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetversions(t *testing.T) {
	gpc := Client{}
	versions, err := gpc.GetVersions(context.Background(), "gopkg.in/src-d/go-billy.v4")
	assert.NoError(t, err)
	assert.Len(t, versions, 11)
}
func TestLatest(t *testing.T) {
	gpc := Client{}
	result, err := gpc.GetLatest(context.Background(), "gopkg.in/src-d/go-billy.v4")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "v4.3.2", result.Version)
}
func TestGetInfo(t *testing.T) {
	version := "v4.3.2"
	gpc := Client{}
	result, err := gpc.GetInfo(context.Background(), "gopkg.in/src-d/go-billy.v4", version)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, version, result.Version)
}
