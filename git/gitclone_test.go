package git

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitClone_should_clone_a_git_repo(t *testing.T) {
	dateisystem, err := GitClone(context.Background(), "https://github.com/halimath/mini-httpd.git", "")
	assert.NoError(t, err)
	datei, err := dateisystem.Open("LICENSE")
	if err != nil {
		t.Fatal(err)
	}
	// defer datei.Close()

	content, err := io.ReadAll(datei)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(content) > 0, "Size of content should be > 0, but was %d", len(content))
	assert.True(t, strings.Contains(string(content), "Apache License"))
	// fmt.Println(string(content))
	// defer os.Remove("LICENSE")

}

func Benchmark_Read_File(b *testing.B) {
	for n := 0; n < b.N; n++ {
		dateisystem, err := GitClone(context.Background(), "https://github.com/halimath/mini-httpd.git", "")
		assert.NoError(b, err)
		datei, err := dateisystem.Open("LICENSE")
		if err != nil {
			b.Fatal(err)
		}
		// defer datei.Close()

		_, err = io.ReadAll(datei)
		if err != nil {
			b.Fatal(err)
		}
	}
}
