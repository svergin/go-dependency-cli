package report

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateReport(t *testing.T) {
	dr := DependencyResolver{}
	// https://github.com/halimath/mini-httpd.git
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	report, err := dr.CreateReport(ctx, "https://github.com/kubernetes/sample-cli-plugin.git", "master")

	if err != nil {
		t.Fatalf("error occured: %v", err)
	}
	assert.Equal(t, 53, len(report.entries))
	report.Print()
}
