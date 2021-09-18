package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bunnier/mdim/internal/base"
	"github.com/bunnier/mdim/internal/cleaner"
	"github.com/bunnier/mdim/internal/markdown"
	"github.com/spf13/cobra"
)

// MdimParam are command-line options.
type MdimParam struct {
	AbsDocFolder string
	AbsImgFolder string

	DoSave           bool
	DoImgDel         bool
	DoWebImgDownload bool
}

var mdimParam = &MdimParam{}

func init() {
	flags := mdimCmd.Flags()
	flags.BoolVarP(&mdimParam.DoSave, "save", "s", false, "Set the option to save markdown document changes, otherwise print scan result only.")
	flags.BoolVarP(&mdimParam.DoImgDel, "delete", "d", false, "Set the option to delete no reference images, otherwise print the paths only.")
	flags.BoolVarP(&mdimParam.DoWebImgDownload, "web", "w", false, "Set the option to download web images to imageFolder. This option might be set with the -s option, otherwise although images have been download to imageFolder, the path in document still be url.")
	flags.StringVarP(&mdimParam.AbsDocFolder, "docFolder", "m", "", "Must not be empty. Assign the folder which markdown documents save in, also can be provided by setting env variable named 'mdim_docFolder'")
	flags.StringVarP(&mdimParam.AbsImgFolder, "imgFolder", "i", "", "Must not be empty. Assign the folder which images save in, also can be provided by setting env variable named 'mdim_imgFolder'.")
}

var mdimCmd = &cobra.Command{
	Use:   "mdim",
	Short: "The tool helps to maintain the images in the markdown files.",
	Long: `The tool helps to maintain the images in the markdown files.
Github: https://github.com/bunnier/mdim`,
	Version: "1.2",
	Run: func(cmd *cobra.Command, args []string) {
		// Try to load folders from env.
		if mdimParam.AbsImgFolder == "" {
			mdimParam.AbsImgFolder = os.Getenv("mdim_imgFolder")
		}

		if mdimParam.AbsDocFolder == "" {
			mdimParam.AbsDocFolder = os.Getenv("mdim_docFolder")
		}

		if mdimParam.AbsImgFolder == "" || mdimParam.AbsDocFolder == "" {
			cmd.Usage()
			os.Exit(1)
		}

		// To abs folder.
		var err error
		if !filepath.IsAbs(mdimParam.AbsImgFolder) {
			mdimParam.AbsImgFolder, err = filepath.Abs(mdimParam.AbsImgFolder)
			if err != nil {
				log.Fatalf("Cannot get the abs path of imageFolder\n%s\n%s", mdimParam.AbsImgFolder, err.Error())
			}
			if _, err := os.Lstat(mdimParam.AbsImgFolder); err != nil {
				log.Fatalf("Cannot get the abs path of imageFolder\n%s\n%s", mdimParam.AbsImgFolder, err.Error())
			}
		}

		if !filepath.IsAbs(mdimParam.AbsDocFolder) {
			mdimParam.AbsDocFolder, err = filepath.Abs(mdimParam.AbsDocFolder)
			if err != nil {
				log.Fatalf("Cannot get the absolutely path of markdownFolder\n%s\n%s", mdimParam.AbsDocFolder, err.Error())
			}
			if _, err := os.Lstat(mdimParam.AbsDocFolder); err != nil {
				log.Fatalf("Cannot get the abs path of markdownFolder\n%s\n%s", mdimParam.AbsDocFolder, err.Error())
			}
		}

		doMdimCmd(mdimParam)
	},
}

func doMdimCmd(param *MdimParam) {
	fmt.Println("==========================================")
	fmt.Println("  Starting to scan markdown document..")
	fmt.Println("==========================================")

	// Scan docs in docFolder to maintain image tags.
	markdownHandleResults := markdown.WalkDirToHandleDocs(
		param.AbsDocFolder,
		param.AbsImgFolder,
		param.DoSave,
		param.DoWebImgDownload)

	hasInterruptErr := false
	allRefImgsAbsPathSet := base.NewSet(100)
	for _, handleResult := range markdownHandleResults {
		if handleResult.HasChangeDuringMaintain ||
			handleResult.RelPathCannotFixedErr != nil ||
			handleResult.WebImgDownloadErr != nil {
			fmt.Println(handleResult.ToString())
			fmt.Println()
		}

		if handleResult.Err != nil {
			hasInterruptErr = true
		}

		allRefImgsAbsPathSet.Merge(handleResult.AllRefImgs)
	}

	if hasInterruptErr {
		os.Exit(10)
	}

	fmt.Println("==========================================")
	fmt.Println("  Starting to scan images..")
	fmt.Println("==========================================")

	// Delete no reference images.
	for _, handleResult := range cleaner.DeleteNoRefImgs(param.AbsImgFolder, allRefImgsAbsPathSet, param.DoImgDel) {
		fmt.Println(handleResult.ToString())
		fmt.Println()
	}

	fmt.Println("==========================================")
	fmt.Println("  All done.")
	fmt.Println("==========================================")
}
