package main

import (
	"fmt"
	"os"

	"github.com/bunnier/mdim/internal/base"
	"github.com/bunnier/mdim/internal/cleaner"
	"github.com/bunnier/mdim/internal/markdown"
	"github.com/spf13/cobra"
)

// MdimOptions are command-line options.
type MdimOptions struct {
	DoImgDel         bool
	DoWebImgDownload bool
}

var mdimOptions = &MdimOptions{}

func init() {
	flags := mdimCmd.Flags()
	initBaseOptions(flags)
	flags.BoolVarP(&mdimOptions.DoImgDel, "delete", "d", false, "Set the option to delete no reference images, otherwise print the paths only.")
	flags.BoolVarP(&mdimOptions.DoWebImgDownload, "web", "w", false, "Set the option to download web images to imageFolder. This option might be set with the '--save' option, otherwise although images have been download to imageFolder, the path in document still be url.")
}

var mdimCmd = &cobra.Command{
	Use:   "mdim",
	Short: "The tool helps to maintain the images in the markdown files.",
	Long: `The tool helps to maintain the images in the markdown files.
Github: https://github.com/bunnier/mdim`,
	Version: "1.2",
	Run: func(cmd *cobra.Command, args []string) {
		validateBaseOptions(cmd)
		doMdimCmd(mdimOptions)
	},
}

func doMdimCmd(param *MdimOptions) {
	fmt.Println("==========================================")
	fmt.Println("  Starting to scan markdown document(s)..")
	fmt.Println("==========================================")

	// workflow steps
	steps := []markdown.ImageMaintainStep{markdown.FixLocalImageRelpathStep}
	if param.DoWebImgDownload {
		steps = append(steps, markdown.DownloadImageStep)
	}

	// Scan docs in docFolder to maintain image tags.
	markdownHandleResults := markdown.WalkDirToHandleDocs(baseOptions.SingleDocument, baseOptions.AbsDocFolder, baseOptions.AbsImgFolder, baseOptions.DoSave, steps)

	hasInterruptErr := false
	allRefImgsAbsPathSet := base.NewSet(100)
	for _, handleResult := range markdownHandleResults {
		if handleResult.HasChangeDuringWorkflow ||
			handleResult.RelPathCannotFixedErr != nil ||
			handleResult.WebImgDownloadErr != nil {
			fmt.Println(handleResult.String())
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
	fmt.Println("  Starting to scan image(s)..")
	fmt.Println("==========================================")

	// Delete no reference images.
	for _, handleResult := range cleaner.DeleteNoRefImgs(baseOptions.AbsImgFolder, allRefImgsAbsPathSet, param.DoImgDel) {
		fmt.Println(handleResult.String())
		fmt.Println()
	}

	fmt.Println("==========================================")
	fmt.Println("  All done.")
	fmt.Println("==========================================")
}
