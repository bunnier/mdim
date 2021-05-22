package main

import (
	"fmt"

	"mdim/core"
)

func main() {
	// Deal with options.
	cliOptions := core.GetOptions()

	fmt.Println("======================================")
	fmt.Println("Starting to scan markdown document..")
	fmt.Println("======================================")

	// Scan docs in docFolder to maintain image tags.
	allRefImgsAbsPathSet, aggErr := core.MaintainImageTags(cliOptions.AbsDocFolder, cliOptions.AbsImgFolder, cliOptions.DoRelPathFix)
	if aggErr != nil {
		aggErr.PrintAggregateError()
		return
	}

	fmt.Println("======================================")
	fmt.Println("Starting to scan images..")
	fmt.Println("======================================")

	// Delete no reference images.
	if errs := core.DeleteNoRefImgs(cliOptions.AbsImgFolder, allRefImgsAbsPathSet, cliOptions.DoImgDel); errs != nil {
		errs.PrintAggregateError()
		return
	}

	fmt.Println("======================================")
	fmt.Println("All done.")
	fmt.Println("======================================")
}
