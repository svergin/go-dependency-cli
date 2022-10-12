package proxyclient

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetversions(t *testing.T) {
	gpc := GoProxyClient{}
	gpc.WithParams("gopkg.in/src-d/go-billy.v4", nil)
	versions, err := gpc.GetVersions(context.Background())
	assert.NoError(t, err)
	assert.Len(t, versions, 11)
}
func TestLatest(t *testing.T) {
	gpc := GoProxyClient{}
	gpc.WithParams("gopkg.in/src-d/go-billy.v4", nil)
	result, err := gpc.GetLatest(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "v4.3.2", result.Version)
}
func TestGetInfo(t *testing.T) {
	version := "v4.3.2"
	gpc := GoProxyClient{}
	gpc.WithParams("gopkg.in/src-d/go-billy.v4", &version)
	result, err := gpc.GetInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, version, result.Version)
}

func TestGetGoMod(t *testing.T) {
	version := "v4.3.2"
	gpc := GoProxyClient{}
	gpc.WithParams("gopkg.in/src-d/go-billy.v4", &version)
	result, err := gpc.GetGoMod(context.Background())
	fmt.Println(result)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, result, "module gopkg.in/src-d/go-billy.v4")
}

func TestErstelleReport(t *testing.T) {
	version := "v4.3.2"
	gpc := GoProxyClient{}
	gpc.WithParams("gopkg.in/src-d/go-billy.v4", &version)
	_, err := gpc.ErstelleReport()
	if err != nil {
		t.Fatal(err)
	}
}
