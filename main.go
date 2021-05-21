package main

import (
	"fmt"

	"mdic/utils"
)

func main() {
	// Deal with options.
	cliOptions := utils.GetOptions()

	fmt.Println("Starting to scan markdown document..")

	// The keys are all image paths.
	allRefImgsMap, errs := utils.ScanToFixImgRelPath(cliOptions.DocFolder, cliOptions.ImgFolder, cliOptions.DoRelPathFix)
	if errs != nil {
		utils.PrintAggregateError(errs)
		return
	}

	fmt.Println("Starting to scan images..")

	// Delete no reference images.
	if errs := utils.DelNoRefImgs(cliOptions.ImgFolder, allRefImgsMap, cliOptions.DoImgDel); errs != nil {
		utils.PrintAggregateError(errs)
		return
	}

	fmt.Println("All done.")
}
