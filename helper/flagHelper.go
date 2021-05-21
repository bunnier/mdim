package helper

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Deal with cli params.
func ParseParams(docFolder *string, imgFolder *string, doFix *bool, doDelete *bool) error {
	var help bool

	flag.Usage = usage
	flag.BoolVar(&help, "h", false, "Show this help.")
	flag.BoolVar(doFix, "f", false, "Set the option to fix image relative paths of markdown documents.")
	flag.BoolVar(doDelete, "d", false, "Set the option to delete no reference images.")
	flag.StringVar(docFolder, "m", "", "The folder markdown documents save in")
	flag.StringVar(imgFolder, "i", "", "The folder images save in")

	flag.Parse()

	if help { // Show usage and then exit directly.
		flag.Usage()
		os.Exit(0)
	}

	switch {
	case *docFolder == "":
		return errors.New("param: missiong -m")
	case *imgFolder == "":
		return errors.New("param: missiong -i")
	}

	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, `mdic - Markdown Images Cleaner

Description: The tool will help you maintain the image relative paths of markdown files and cleanup no reference images.

Github: https://github.com/bunnier/mdic

Usage: mdic [-dfh] [-i image folder] [-m markdown doc folder] 

Options:
`)
	flag.PrintDefaults()
}
