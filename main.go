package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "gomove"
	app.Usage = "Move Golang packages to a new path."
	app.Version = "0.2.17"
	app.ArgsUsage = "[old path] [new path]"
	app.Author = "Kaushal Subedi <kaushal@subedi.co>"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dir, d",
			Value: "./",
			Usage: "directory to scan",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "only move imports in a file",
		},
		cli.StringFlag{
			Name:  "prefix,p",
			Value: "true",
			Usage: "path prefix will be used(true)     path exact matched(false)",
		},
		cli.StringFlag{
			Name:  "safe-mode, s",
			Value: "false",
			Usage: "run program in safe mode (comments will be wiped)",
		},
	}

	app.Action = func(c *cli.Context) {
		file := c.String("file")
		dir := c.String("dir")
		usePreifx := c.String("prefix")
		from := c.Args().Get(0)
		to := c.Args().Get(1)

		if file != "" {
			ProcessFile(file, from, to, usePrefix, c)
		} else {
			ScanDir(dir, from, to, usePrefix, c)
		}

	}

	app.Run(os.Args)
}

// ScanDir scans a directory for go files and
func ScanDir(dir string, from string, to string, usePrefix string, c *cli.Context) {

	// If from and to are not empty scan all files
	if from != "" && to != "" {
		// Scan directory for files
		filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
			if info.IsDir() && info.Name() == "vendor" {
				return filepath.SkipDir
			}
			// Only process go files
			if path.Ext(filePath) == ".go" {
				ProcessFile(filePath, from, to, usePrefix, c)
			}
			return nil
		})

	} else {
		cli.ShowAppHelp(c)
	}

}

// ProcessFile processes file according to what mode is chosen
func ProcessFile(filePath string, from string, to string, prefixMode string, c *cli.Context) {
	usePrefix := prefixMode == "true"
	if c.String("safe-mode") == "true" {
		ProcessFileAST(filePath, from, to, usePrefix)
	} else {
		ProcessFileNative(filePath, from, to, usePrefix)
	}
}
