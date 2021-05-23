package main

import (
	"fmt"
	"os"

	"mdim/core"
)

func main() {
	// Deal with options.
	cliOptions := core.GetOptions()

	fmt.Println("==========================================")
	fmt.Println("  Starting to scan markdown document..")
	fmt.Println("==========================================")

	// Scan docs in docFolder to maintain image tags.
	allRefImgsAbsPathSet, markdownHandleResults := core.MaintainImageTags(cliOptions.AbsDocFolder, cliOptions.AbsImgFolder, cliOptions.DoRelPathFix)
	hasInteruptErr := false
	for _, handleResult := range markdownHandleResults {
		if handleResult.HasErrImgRelPath || handleResult.RelPathCannotFixedErr != nil {
			fmt.Println(handleResult.ToString())
			fmt.Println()
		}
		if handleResult.Err != nil {
			hasInteruptErr = true
		}
	}

	if hasInteruptErr {
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
