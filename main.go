package main

import (
	"fmt"

	"mdic/utils"
)

func main() {
	// Deal with options.
	var cliOptions *utils.CliOptions
	var err error
	if cliOptions, err = utils.GetOptions(); err != nil {
		println(err.Error())
		return
	}

	fmt.Println("Starting to scan markdown document..")

	// The keys are all image paths.
	allRefImagesMap, errs := utils.ScanToFixImgRelPath(cliOptions.DocFolder, cliOptions.ImgFolder, cliOptions.DoFix)
	if errs != nil {
		utils.PrintAggregateError(errs)
		return
	}

	fmt.Println("Starting to scan images..")

	// Delete no reference images.
	if errs := utils.DelNoRefImags(cliOptions.ImgFolder, allRefImagesMap, cliOptions.DoDel); errs != nil {
		utils.PrintAggregateError(errs)
		return
	}

	fmt.Println("All done.")
}
