package main

import (
	"io"
	"os"
)

type ByAscending []string

func (a ByAscending) Len() int {
	return len(a)
}

func (a ByAscending) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByAscending) Less(i, j int) bool {
	return a[i] < a[j]
}

func dirTree(out io.Writer, path string, printFiles bool) error {

	return nil
}

func main() {
	out := os.Stdout

	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	isPrintFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, isPrintFiles)
	if err != nil {
		panic(err.Error())
	}
}
