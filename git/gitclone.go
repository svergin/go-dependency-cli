package git

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func GitClone(ctx context.Context, url, branch string) (fs.FS, error) {
	if branch == "" {
		branch = "main"
	}
	opts := &gogit.CloneOptions{
		URL:           url,
		Progress:      os.Stdout,
		RemoteName:    branch,
		SingleBranch:  true,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
	}

	billyFs := memfs.New()
	_, err := gogit.Clone(memory.NewStorage(), billyFs, opts)
	if err != nil {
		return nil, err
	}

	myFS := InMemFS{
		bfs: billyFs,
	}

	return &myFS, nil
}

type InMemFS struct {
	bfs billy.Filesystem
}

type InMemFile struct {
	bf  billy.File
	bfs billy.Filesystem
}

func (imf InMemFile) Stat() (fs.FileInfo, error) {
	return imf.bfs.Stat(imf.bf.Name())
}

func (imf InMemFile) Read(b []byte) (int, error) {
	return imf.bf.Read(b)
}

func (imf InMemFile) Close() error {
	return imf.bf.Close()
}

func (imfs InMemFS) Open(name string) (fs.File, error) {
	bfile, err := imfs.bfs.Open(name)
	if err != nil {
		return nil, err
	}
	imf := InMemFile{
		bf:  bfile,
		bfs: imfs.bfs,
	}
	return imf, nil
}

func (imfs InMemFS) Open_(name string) (fs.File, error) {
	bfile, err := imfs.bfs.Open(name)
	if err != nil {
		return nil, err
	}
	defer bfile.Close()
	myfile, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(myfile, bfile)
	if err != nil {
		return nil, err
	}
	myfile.Seek(0, 0)
	return myfile, nil

}

func (imfs InMemFS) Open__(name string) (fs.File, error) {

	r := imfs.bfs.Root()

	files, err := imfs.bfs.ReadDir(r)
	if err != nil {
		return nil, err
	}

	return imfs.findFile(name, files)
}

func (imfs InMemFS) findFile(name string, files []fs.FileInfo) (fs.File, error) {
	var fileToReturn fs.File
	for _, file := range files {
		if file.IsDir() {
			moreFiles, err := imfs.bfs.ReadDir(file.Name())
			if err != nil {
				return nil, err
			}
			imfs.findFile(name, moreFiles)
		} else {
			if file.Name() == name {
				bf, err := imfs.bfs.OpenFile(file.Name(), os.O_RDONLY, file.Mode().Perm())
				if err != nil {
					return nil, err
				}
				fileToReturn, err := os.Create(name)
				if err != nil {
					return nil, err
				}
				defer fileToReturn.Seek(0, 0)
				buf := make([]byte, 1024)
				for {
					n, err := bf.Read(buf)

					if err == io.EOF || n == 0 {
						break
					}
					if n > 0 {
						_, err = fileToReturn.Write(buf[:n])
						if err != nil {
							return nil, err
						}
					}
				}
				defer bf.Close()
				fileToReturn.Sync()
				return fileToReturn, nil
			} else {
				continue
			}
		}

	}
	return fileToReturn, nil
}
