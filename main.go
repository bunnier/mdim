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
	for _, handleResult := range core.DeleteNoRefImgs(cliOptions.AbsImgFolder, allRefImgsAbsPathSet, cliOptions.DoImgDel) {
		switch {
		case handleResult.Err != nil:
			fmt.Printf("[image handle]:find a no reference image, but fail to delete.\n----> %s\n", handleResult.ImagePath)
		case handleResult.Deleted:
			fmt.Printf("[image handle]:delete a no reference image successfully.\n----> %s\n", handleResult.ImagePath)
		case !handleResult.Deleted:
			fmt.Printf("[image handle]:find a no reference image, do not delete this time.\n----> %s\n", handleResult.ImagePath)
		}
	}

	fmt.Println("======================================")
	fmt.Println("All done.")
	fmt.Println("======================================")
}
