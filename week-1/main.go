package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

const pathSeparator = string(os.PathSeparator)

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

type finder func(path string) ([]string, error)

func dirs(path string) ([]string, error) {
	var dirs []string

	f, err := os.Open(path)
	if err != nil {
		return dirs, err
	}

	tree, err := f.Readdir(-1)
	if err != nil {
		return dirs, err
	}

	for _, file := range tree {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}

	return dirs, nil
}

func files(path string) ([]string, error) {
	var files []string

	f, err := os.Open(path)
	if err != nil {
		return files, err
	}

	tree, err := f.Readdir(-1)
	if err != nil {
		return files, err
	}

	for _, file := range tree {
		if !file.IsDir() {
			content := fileInfoContent(file)
			files = append(files, content)
		} else {
			files = append(files, file.Name())
		}
	}

	return files, nil
}

func fileInfoContent(file os.FileInfo) string {
	var info string

	if file.Size() > 0 {
		info = fmt.Sprintf("%s (%db)", file.Name(), file.Size())
	} else {
		info = fmt.Sprintf("%s (empty)", file.Name())
	}

	return info
}

func buildTree(path string, prefix string, f finder) (string, error) {
	var files, err = f(path)
	var tree string

	if len(files) == 0 {
		return tree, nil
	}

	if err != nil {
		return tree, err
	}

	sort.Sort(ByAscending(files))

	lastIdx := len(files) - 1

	for _, file := range files[:lastIdx] {
		subPath := fmt.Sprint(path, pathSeparator, file)
		subPrefix := fmt.Sprint(prefix, "│	")

		subtrees, err := buildTree(subPath, subPrefix, f)
		if err != nil {
			return tree, err
		}

		tree = fmt.Sprint(tree, prefix, "├───", file, "\n", subtrees)
	}

	subPath := fmt.Sprint(path, pathSeparator, files[lastIdx])
	subPrefix := fmt.Sprint(prefix, "	")

	subtrees, err := buildTree(subPath, subPrefix, f)
	if err != nil {
		return tree, err
	}

	return fmt.Sprint(tree, prefix, "└───", files[lastIdx], "\n", subtrees), nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	f := dirs
	if printFiles {
		f = files
	}

	tree, err := buildTree(path, "", f)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(out, tree)

	return err
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
