package utils

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Command-line options
type CliOptions struct {
	DocFolder string
	ImgFolder string
	DoFix     bool
	DoDel     bool
}

// Deal with cli params.
func GetOptions() (*CliOptions, error) {
	var help bool
	CliParams := &CliOptions{}

	flag.Usage = usage
	flag.BoolVar(&help, "h", false, "Show this help.")
	flag.BoolVar(&CliParams.DoFix, "f", false, "Set the option to fix image relative paths of markdown documents.")
	flag.BoolVar(&CliParams.DoDel, "d", false, "Set the option to delete no reference images.")
	flag.StringVar(&CliParams.DocFolder, "m", "", "The folder markdown documents save in")
	flag.StringVar(&CliParams.ImgFolder, "i", "", "The folder images save in")

	flag.Parse()

	if help { // Show usage and then exit directly.
		flag.Usage()
		os.Exit(0)
	}

	switch {
	case CliParams.DocFolder == "":
		return nil, errors.New("param: missiong -m")
	case CliParams.ImgFolder == "":
		return nil, errors.New("param: missiong -i")
	}

	return CliParams, nil
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
