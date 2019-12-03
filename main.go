package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Writer interface {
	WriteString(str string) (n int, err error)
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

func dirTree(out Writer, path string, printFiles bool) error {
	writeDirTree(out, path, printFiles, "")
	return nil
}

func writeDirTree(out Writer, path string, printFiles bool, prefix string) error {
	dir, err := os.Open(path)
	if err != nil {
		fmt.Errorf("error to Open dir")
	}

	dirNames, _ := dir.Readdir(-1)
	dir.Close()

	if printFiles == false {
		dirNames = getOnlyFolders(dirNames)
	}

	sort.Slice(dirNames, func(i, j int) bool { return dirNames[i].Name() < dirNames[j].Name() })
	for index, element := range dirNames {
		if element.Name() == ".DS_Store" {
			continue
		}

		isLast := index == len(dirNames)-1
		writeElement(out, element.Name(), element.Size(), element.IsDir(), isLast, prefix)
		if element.IsDir() {
			newPath := string(path + string(os.PathSeparator) + element.Name())
			writeDirTree(out, newPath, printFiles, getPrefixForNextLevel(prefix, isLast))
		}
	}

	return nil
}

func getOnlyFolders(slise []os.FileInfo) []os.FileInfo {
	result := []os.FileInfo{}
	for _, element := range slise {
		if element.IsDir() {
			result = append(result, element)
		}
	}

	return result
}

func writeElement(out Writer, name string, size int64, isDir bool, isLast bool, prefix string) {
	sizeText := ""
	if !isDir {
		if size > 0 {
			sizeText = " (" + strconv.FormatInt(size, 10) + "b)"
		} else {
			sizeText = " (empty)"
		}
	}

	if isLast {
		prefix += "└───"
	} else {
		prefix += "├───"
	}

	out.WriteString(prefix + name + sizeText + "\n")

}

func getPrefixForNextLevel(oldPrefix string, isLastElement bool) string {
	if isLastElement {
		return oldPrefix + "\t"
	}

	return oldPrefix + "│\t"
}
