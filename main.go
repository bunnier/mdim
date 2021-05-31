package main

import (
	"fmt"
	"os"

	"1024baby.com/mdim/core"
)

func main() {
	// Deal with options.
	cliOptions := core.GetOptions()

	fmt.Println("==========================================")
	fmt.Println("  Starting to scan markdown document..")
	fmt.Println("==========================================")

	// Scan docs in docFolder to maintain image tags.
	allRefImgsAbsPathSet, markdownHandleResults := core.MaintainImageTags(
		cliOptions.AbsDocFolder,
		cliOptions.AbsImgFolder,
		cliOptions.DoSave,
		cliOptions.DoWebImgDownload)

	hasInterruptErr := false
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
	}

	if hasInterruptErr {
		os.Exit(10)
	}

	fmt.Println("==========================================")
	fmt.Println("  Starting to scan images..")
	fmt.Println("==========================================")

	// Delete no reference images.
	for _, handleResult := range core.DeleteNoRefImgs(cliOptions.AbsImgFolder, allRefImgsAbsPathSet, cliOptions.DoImgDel) {
		fmt.Println(handleResult.ToString())
		fmt.Println()
	}

	fmt.Println("==========================================")
	fmt.Println("  All done.")
	fmt.Println("==========================================")
}
