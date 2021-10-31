// Package main generate documentation for plugins.
package main

import (
	"fmt"
	"go/parser"
	"log"
	"os"
	"path/filepath"

	"github.com/josephspurrier/ambient/cmd/docgen/lib"
)

func main() {
	dir := "plugin/prism"
	fmt.Println("Folder:", dir)

	// Go up two folders.
	root, err := filepath.Abs("../..")
	if err != nil {
		log.Fatalln(err.Error())
	}

	absDir := filepath.Join(root, dir)

	// Ensure the folder exists.
	_, err = os.ReadDir(absDir)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Get the folder name
	dirName := filepath.Base(dir)
	fmt.Println(dirName)

	fileName := filepath.Join(absDir, dirName+".go")

	_, err = os.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err.Error())
	}

	//fmt.Println(string(mainFile))

	// Read the file
	gt, err := lib.Load(fileName, parser.ParseComments)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(gt.PrintComments())
}
