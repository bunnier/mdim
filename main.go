package main

import (
	"fmt"

	"mdic/helper"
)

func main() {
	// Deal with options.
	var cliOptions *helper.CliOptions
	var err error
	if cliOptions, err = helper.GetOptions(); err != nil {
		println(err.Error())
		return
	}

	fmt.Println("Starting to scan markdown document..")

	// The keys are all image paths.
	allRefImagesMap, errs := helper.ScanToFixImgRelPath(cliOptions.DocFolder, cliOptions.ImgFolder, cliOptions.DoFix)
	if errs != nil {
		helper.PrintAggregateError(errs)
		return
	}

	fmt.Println("Starting to scan images..")

	// Delete no reference images.
	if errs := helper.DelNoRefImags(cliOptions.ImgFolder, allRefImagesMap, cliOptions.DoDel); errs != nil {
		helper.PrintAggregateError(errs)
		return
	}

	fmt.Println("All done.")
}
