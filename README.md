# go-dependency-cli
# Session 1.1 - Build a VCS client
Build an abstraction package that clones a VCS repo and provides access to the cloned files via 
[`fs.FS`](https://pkg.go.dev/io/fs#FS). The abstraction should basically be a function with a signature like
```go
func GitClone (ctx context.Context, url, branch string) (fs.FS, error)
```
* [github.com/go-git/go-git](https://github.com/go-git/go-git) will be very helpful to implement the operation
* `go-git` uses a [`billy.Filesystem`](https://pkg.go.dev/github.com/go-git/go-billy/v5#Filesystem)
 which does not diretly matches `fs.FS`. You need to work a way around.
* You can do the clone "in-memory" or into a temporary directory. What is the performance difference?
* Write a unit test. You can use https://github.com/halimath/mini-httpd.git (which is very small and fast to
 clone) and assert that the file `LICENSE` has a first line with trimmed text equal to `"Apache License"`.
