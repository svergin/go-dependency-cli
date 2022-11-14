package report

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	dr := DependencyResolver{}

	report, err := dr.ErstelleReport(context.Background(), "https://github.com/halimath/mini-httpd.git", "main")
	if err != nil {
		t.Fatalf("error occured: %v", err)
	}
	assert.Equal(t, 1, len(report.entries))
}
