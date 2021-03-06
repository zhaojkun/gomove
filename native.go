package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mgutz/ansi"
)

// ProcessFileNative processes files uing native string search instead of ast parsing
func ProcessFileNative(filePath string, from string, to string, usePrefix bool) {
	if !usePrefix {
		from = fmt.Sprintf(`"%s"`, from)
		to = fmt.Sprintf(`"%s"`, to)
	}
	//Colors to be used on the console
	red := ansi.ColorCode("red+bh")
	white := ansi.ColorCode("white+bh")
	green := ansi.ColorCode("green+bh")
	yellow := ansi.ColorCode("yellow+bh")
	blackOnWhite := ansi.ColorCode("black+b:white+h")
	//Reset the color
	reset := ansi.ColorCode("reset")

	fmt.Println(blackOnWhite+"Processing file", filePath, reset)

	// Open file to read
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}

	// Scan file line by line
	scanner := bufio.NewScanner(bytes.NewReader(fileContent))

	// Track line that is being scanned
	scanLine := 0
	// Track number of changes in file
	numChanges := 0

	// Control variables
	isImportLine := false

	// Store final output text
	output := ""

	// Scan through the lines of go file
	for scanner.Scan() {

		scanLine++
		line := scanner.Text()
		bareLine := strings.Replace(line, " ", "", -1)

		// If it is a single import statement, replace the path in that line
		if strings.Contains(bareLine, "import\"") {
			newImport := strings.Replace(line, from, to, -1)
			output += newImport + "\n"
			if line != newImport {
				numChanges++

				fmt.Println(red+"Updating "+
					reset+white+
					strings.TrimSpace(line)+
					reset+red+
					" to "+
					reset+white+
					strings.TrimSpace(newImport)+
					reset+red+
					"on line", scanLine, reset)
			}

			continue
		}

		// Change isImportLine accordingly if import statements are detected
		if strings.Contains(bareLine, "import(") {
			fmt.Println(green+"Found Multiple Imports Starting On Line", scanLine, reset)
			isImportLine = true
		} else if isImportLine && strings.Contains(bareLine, ")") {
			fmt.Println(green+"Imports Finish On Line", scanLine, reset)
			isImportLine = false
		}

		// If it is a import line, replace the import
		if isImportLine {
			newImport := strings.Replace(line, from, to, -1)
			output += newImport + "\n"
			if line != newImport {
				numChanges++
				fmt.Println(red+"Updating text "+
					reset+white+
					strings.TrimSpace(line)+
					reset+red+
					" to "+
					reset+white+
					strings.TrimSpace(newImport)+
					reset+red+
					" on line", scanLine, reset)

			}
			continue
		}

		// Just copy the rest of the lines to the output
		output += line + "\n"

	}

	// Only write if changes were made
	if numChanges > 0 {
		fmt.Println(yellow+
			"File",
			filePath,
			"saved after",
			numChanges,
			"changes",
			reset, "\n\n")
		formatedSource, err := format.Source([]byte(output))
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(filePath, formatedSource, os.ModePerm)
	} else {
		fmt.Println(yellow+
			"No changes to write on this file.",
			reset, "\n\n")
	}
}
