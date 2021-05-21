package core

import (
	"flag"
	"fmt"
	"os"
)

// Command-line options.
type CliOptions struct {
	DocFolder    string
	ImgFolder    string
	DoRelPathFix bool
	DoImgDel     bool
}

// Deal with cli params.
func GetOptions() *CliOptions {
	var help bool
	CliParams := &CliOptions{}

	flag.Usage = usage
	flag.BoolVar(&help, "h", false, "Show this help.")
	flag.BoolVar(&CliParams.DoRelPathFix, "f", false, "Set the option to fix image relative paths of markdown documents.")
	flag.BoolVar(&CliParams.DoImgDel, "d", false, "Set the option to delete no reference images.")
	flag.StringVar(&CliParams.DocFolder, "m", "", "Must be not empty. The folder markdown documents save in")
	flag.StringVar(&CliParams.ImgFolder, "i", "", "Must be not empty. The folder images save in")

	flag.Parse()

	// Show usage and then exit directly.
	if help || CliParams.ImgFolder == "" || CliParams.DocFolder == "" {
		flag.Usage()
		os.Exit(0)
	}

	return CliParams
}

func usage() {
	fmt.Fprintf(os.Stderr, `mdic - Markdown Images Cleaner

Description: The tool will help you maintain the image relative paths of markdown files and cleanup no reference images.

Github: https://github.com/bunnier/mdic

Usage: mdic [-h] [-d] [-f] [-i imageFolder] [-m markdownFolder] 

Options:
`)
	flag.PrintDefaults()
}
