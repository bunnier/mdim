package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	mdimCmd.AddCommand(qiniuCmd)
	mdimCmd.Execute()
}

type BaseOptions struct {
	AbsDocFolder   string
	AbsImgFolder   string
	SingleDocument string
	DoSave         bool
}

var baseOptions = &BaseOptions{}

func initBaseOptions(flags *pflag.FlagSet) {
	flags.StringVarP(&baseOptions.SingleDocument, "doc", "m", "", "Assign the target markdown document. There must be at least one of '--doc' and '--docFolder'.")
	flags.StringVarP(&baseOptions.AbsDocFolder, "docFolder", "f", "", "Assign the folder which markdown documents save in, also can be provided by setting env variable named 'mdim_docFolder'")
	flags.StringVarP(&baseOptions.AbsImgFolder, "imgFolder", "i", "", "Must not be empty. Assign the folder which images save in, also can be provided by setting env variable named 'mdim_imgFolder'.")
	flags.BoolVarP(&baseOptions.DoSave, "save", "s", false, "Set the option to save markdown document changes, otherwise print scan result only.")
}

func validateBaseOptions(cmd *cobra.Command) {
	// Try to load folders from env.
	if baseOptions.AbsImgFolder == "" {
		baseOptions.AbsImgFolder = os.Getenv("mdim_imgFolder")
	}

	if baseOptions.AbsDocFolder == "" {
		baseOptions.AbsDocFolder = os.Getenv("mdim_docFolder")
	}

	if baseOptions.AbsImgFolder == "" || (baseOptions.AbsDocFolder == "" && baseOptions.SingleDocument == "") {
		cmd.Usage()
		os.Exit(1)
	}

	// To abs folder.
	var err error
	if !filepath.IsAbs(baseOptions.AbsImgFolder) {
		baseOptions.AbsImgFolder, err = filepath.Abs(baseOptions.AbsImgFolder)
		if err != nil {
			log.Fatalf("Cannot get the abs path of imageFolder.\n%s\n%s", baseOptions.AbsImgFolder, err.Error())
		}
		if _, err := os.Lstat(baseOptions.AbsImgFolder); err != nil {
			log.Fatalf("Cannot get the abs path of imageFolder.\n%s\n%s", baseOptions.AbsImgFolder, err.Error())
		}
	}

	if baseOptions.SingleDocument != "" {
		if !filepath.IsAbs(baseOptions.SingleDocument) {
			baseOptions.SingleDocument, err = filepath.Abs(baseOptions.SingleDocument)
			if err != nil {
				log.Fatalf("Cannot get the absolutely path of document.\n%s\n%s", baseOptions.SingleDocument, err.Error())
			}
		}

		if _, err := os.Stat(baseOptions.SingleDocument); err != nil {
			log.Fatalf("Please check the document path.\n%s\n%s", baseOptions.SingleDocument, err.Error())
		}
	} else {
		if !filepath.IsAbs(baseOptions.AbsDocFolder) {
			baseOptions.AbsDocFolder, err = filepath.Abs(baseOptions.AbsDocFolder)
			if err != nil {
				log.Fatalf("Cannot get the absolutely path of markdownFolder.\n%s\n%s", baseOptions.AbsDocFolder, err.Error())
			}
			if _, err := os.Lstat(baseOptions.AbsDocFolder); err != nil {
				log.Fatalf("Cannot get the abs path of markdownFolder.\n%s\n%s", baseOptions.AbsDocFolder, err.Error())
			}
		}
	}
}
