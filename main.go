package main

import (
	"context"
	"fmt"
	"go-dependency-cli/report"
	"os"
	"time"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Falsche Anzahl an Argumenten.")
		fmt.Println("Beispiele:")
		fmt.Println("Checkout Main-Branch: 'go run ./main.go https://github.com/kubernetes/sample-cli-plugin.git'")
		fmt.Println("Checkout bestimmten Branch: 'go run ./main.go https://github.com/kubernetes/sample-cli-plugin.git master'")
		os.Exit(-1)
	}
	repoUrl := args[1]
	var branch string
	if len(args) == 3 {
		branch = args[2]
	}
	resolver := report.DependencyResolver{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	report, err := resolver.CreateReport(ctx, repoUrl, branch)
	if err != nil {
		fmt.Printf("error ocurred: %v\n", err)
		os.Exit(-1)
	}
	report.Print()
}
