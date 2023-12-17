package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func printSize(out io.Writer, dirEntry os.DirEntry) {
	info, _ := dirEntry.Info()

	out.Write([]byte(" ("))
	if info.Size() == 0 {
		out.Write([]byte("empty"))
	} else {
		out.Write([]byte(strconv.Itoa(int(info.Size())) + "b"))
	}
	out.Write([]byte(")"))
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := draw(out, path, "", printFiles)
	if err != nil {
		return fmt.Errorf("Error while drawing the graph: %s\n", err)
	}

	return nil
}

func lastDir(dirEntries []os.DirEntry) int {
	lastDirIndex := len(dirEntries) - 1

	for i := len(dirEntries) - 1; i >= 0; i-- {
		if dirEntries[i].IsDir() {
			lastDirIndex = i
			break
		}
	}

	return lastDirIndex
}

func draw(out io.Writer, path string, prefix string, printFiles bool) error {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("Error while reading the directory path: %s", err)
	}

	lastDirIndex := len(dirEntries) - 1
	if !printFiles {
		lastDirIndex = lastDir(dirEntries)
	}

	for index, value := range dirEntries {
		if !printFiles && !value.IsDir() {
			continue
		}

		out.Write([]byte(prefix))

		descPrefix := prefix
		if index == len(dirEntries)-1 || (!printFiles && index == lastDirIndex) {
			out.Write([]byte("└───" + value.Name()))
			descPrefix += "\t"
		} else {
			out.Write([]byte("├───" + value.Name()))
			descPrefix += "│\t"
		}

		if printFiles && !value.IsDir() {
			printSize(out, value)
		}

		out.Write([]byte("\n"))

		draw(out, path+"/"+value.Name(), descPrefix, printFiles)
	}

	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
