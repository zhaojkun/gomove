package main

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mgutz/ansi"
	"golang.org/x/tools/go/ast/astutil"
)

// ProcessFileAST processes the files using golang's AST parser
func ProcessFileAST(filePath string, from string, to string, usePrefix bool) {

	//Colors to be used on the console
	red := ansi.ColorCode("red+bh")
	white := ansi.ColorCode("white+bh")
	yellow := ansi.ColorCode("yellow+bh")
	blackOnWhite := ansi.ColorCode("black+b:white+h")
	//Reset the color
	reset := ansi.ColorCode("reset")

	fmt.Println(blackOnWhite+"Processing file", filePath, "in SAFE MODE", reset)

	// New FileSet to parse the go file to
	fSet := token.NewFileSet()

	// Parse the file
	file, err := parser.ParseFile(fSet, filePath, nil, 0)
	if err != nil {
		fmt.Println(err)
	}

	// Get the list of imports from the ast
	imports := astutil.Imports(fSet, file)

	// Keep track of number of changes
	numChanges := 0

	// Iterate through the imports array
	for _, mPackage := range imports {
		for _, mImport := range mPackage {
			// Since astutil returns the path string with quotes, remove those
			importString := strings.TrimSuffix(strings.TrimPrefix(mImport.Path.Value, "\""), "\"")
			var matched bool
			if usePrefix {
				matched = strings.Contains(importString, from)
			} else {
				matched = importString == from
			}
			// If the path matches the oldpath, replace it with the new one
			if matched {
				//If it needs to be replaced, increase numChanges so we can write the file later
				numChanges++

				// Join the path of the import package with the remainder from the old one after removing the old import package
				replacePackage := strings.Replace(importString, from, to, -1)

				fmt.Println(red +
					"Updating import " +
					reset + white +
					importString +
					reset + red +
					" to " +
					reset + white +
					replacePackage +
					reset)

				// Remove the old import and replace it with the replacement
				astutil.DeleteImport(fSet, file, importString)
				astutil.AddImport(fSet, file, replacePackage)
			}
		}
	}

	// If the number of changes are more than 0, write file
	if numChanges > 0 {
		// Print the new AST tree to a new output buffer
		var outputBuffer bytes.Buffer
		printer.Fprint(&outputBuffer, fSet, file)
		output, err := format.Source(outputBuffer.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(filePath, output, os.ModePerm)
		fmt.Println(yellow+
			"File",
			filePath,
			"saved after",
			numChanges,
			"changes",
			reset, "\n\n")
	} else {
		fmt.Println(yellow+
			"No changes to write on this file.",
			reset, "\n\n")
	}
}
