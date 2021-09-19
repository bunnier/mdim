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
	RelFix           bool
	DoImgDel         bool
	DoWebImgDownload bool
}

var mdimOptions = &MdimOptions{}

func init() {
	flags := mdimCmd.Flags()
	initBaseOptions(flags)
	flags.BoolVarP(&mdimOptions.RelFix, "relfix", "r", false, "Set this option to fix wrong local relative path of images after moved document.")
	flags.BoolVarP(&mdimOptions.DoImgDel, "delete", "d", false, "Set this option to delete no reference images, otherwise print the paths only.")
	flags.BoolVarP(&mdimOptions.DoWebImgDownload, "web", "w", false, "Set this option to download reference web images to 'imgFolder'.")
}

var mdimCmd = &cobra.Command{
	Use:   "mdim",
	Short: "The tool helps to maintain the images in the markdown files.",
	Long: `The tool helps to maintain the images in the markdown files.
Github: https://github.com/bunnier/mdim`,
	Version: "2.0.0",
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
	steps := make([]markdown.ImageMaintainStep, 0, 2)

	if param.RelFix {
		steps = append(steps, markdown.FixLocalImageRelpathStep)
	}

	if param.DoWebImgDownload {
		steps = append(steps, markdown.DownloadImageStep)
	}

	// Scan docs in docFolder to maintain image tags.
	markdownHandleResults := markdown.WalkDirToHandleDocs(baseOptions.SingleDocument, baseOptions.AbsDocFolder, baseOptions.AbsImgFolder, steps)

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

	if param.DoImgDel {
		fmt.Println("==========================================")
		fmt.Println("  Starting to scan image(s)..")
		fmt.Println("==========================================")

		// Delete no reference images.
		for _, handleResult := range cleaner.DeleteNoRefImgs(baseOptions.AbsImgFolder, allRefImgsAbsPathSet) {
			fmt.Println(handleResult.String())
			fmt.Println()
		}
	}

	fmt.Println("==========================================")
	fmt.Println("  All done.")
	fmt.Println("==========================================")
}
