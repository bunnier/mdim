package core

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// CliOptions are command-line options.
type CliOptions struct {
	AbsDocFolder     string
	AbsImgFolder     string
	DoSave           bool
	DoImgDel         bool
	DoWebImgDownload bool
}

// GetOptions from command-line options.
func GetOptions() *CliOptions {
	var help bool
	CliParams := &CliOptions{}

	flag.Usage = usage
	flag.BoolVar(&help, "h", false, "Show this help.")
	flag.BoolVar(&CliParams.DoSave, "s", false, "Set the option to save markdown document changes, otherwise print scan result only.")
	flag.BoolVar(&CliParams.DoImgDel, "d", false, "Set the option to delete no reference images, otherwise print the paths only.")
	flag.BoolVar(&CliParams.DoWebImgDownload, "w", false, "Set the option to download web images to imageFolder. This option often be set with the -s option, otherwise although images have been download to imageFolder, the path in document still be web paths.")
	flag.StringVar(&CliParams.AbsDocFolder, "m", "", "Must not be empty. Assign the folder which markdown documents save in, also can be provided by setting env variable named 'mdim_docFolder'")
	flag.StringVar(&CliParams.AbsImgFolder, "i", "", "Must not be empty. Assign the folder which images save in, also can be provided by setting env variable named 'mdim_imgFolder'.")

	flag.Parse()

	// Show usage and then exit directly.
	if help {
		flag.Usage()
		os.Exit(0)
	}

	if CliParams.AbsImgFolder == "" {
		CliParams.AbsImgFolder = os.Getenv("mdim_imgFolder")
	}

	if CliParams.AbsDocFolder == "" {
		CliParams.AbsDocFolder = os.Getenv("mdim_docFolder")
	}

	if CliParams.AbsImgFolder == "" || CliParams.AbsDocFolder == "" {
		flag.Usage()
		os.Exit(1)
	}

	var err error
	if !filepath.IsAbs(CliParams.AbsImgFolder) {
		CliParams.AbsImgFolder, err = filepath.Abs(CliParams.AbsImgFolder)
		if err != nil {
			fmt.Printf("Cannot get the abs path of imageFolder\n%s\n%s", CliParams.AbsImgFolder, err.Error())
			os.Exit(2)
		}
		if _, err := os.Lstat(CliParams.AbsImgFolder); err != nil {
			fmt.Printf("Cannot get the abs path of imageFolder\n%s\n%s", CliParams.AbsImgFolder, err.Error())
			os.Exit(3)
		}
	}

	if !filepath.IsAbs(CliParams.AbsDocFolder) {
		CliParams.AbsDocFolder, err = filepath.Abs(CliParams.AbsDocFolder)
		if err != nil {
			fmt.Printf("Cannot get the absolutely path of markdownFolder\n%s\n%s", CliParams.AbsDocFolder, err.Error())
			os.Exit(4)
		}
		if _, err := os.Lstat(CliParams.AbsDocFolder); err != nil {
			fmt.Printf("Cannot get the abs path of markdownFolder\n%s\n%s", CliParams.AbsDocFolder, err.Error())
			os.Exit(5)
		}
	}

	return CliParams
}

func usage() {
	fmt.Fprintf(os.Stderr, `mdim - Markdown Images Maintainer

Description: The tool will help to maintain the images in markdown files.

Github: https://github.com/bunnier/mdim

Usage: mdim [-h] [-d] [-w] [-i imageFolder] [-m markdownFolder] 

Options:
`)
	flag.PrintDefaults()
}
