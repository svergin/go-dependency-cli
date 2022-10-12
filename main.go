package main

import (
	"fmt"
	"go-dependency-cli/proxyclient"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Falsche Anzahl an Argumenten. Beispiel-Aufruf: 'go run ./main.go gopkg.in/src-d/go-billy.v4' ")
	}
	modulename := args[1]
	gpc := proxyclient.GoProxyClient{}
	gpc.WithParams(modulename, nil)
	gpc.ErstelleReport()
}
