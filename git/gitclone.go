package git

import (
	"context"
	"io/fs"
	"os"

	git "gopkg.in/src-d/go-git.v4"
)

func GitClone(ctx context.Context, url, branch string) (fs.FS, error) {
	if branch == "" {
		branch = "master"
	}
	opts := &git.CloneOptions{
		URL:        url,
		Progress:   os.Stdout,
		RemoteName: branch,
	}
	repo, err := git.PlainClone("/tmp/repo", false, opts)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
