package main

import (
	"fmt"

	"mdic/helper"
)

func main() {
	var (
		docFolder string
		imgFolder string
		doFix     bool
		doDel     bool
	)

	// Deal with options.
	if err := helper.ParseParams(&docFolder, &imgFolder, &doFix, &doDel); err != nil {
		println(err.Error())
		return
	}

	fmt.Println("Starting to scan markdown document..")

	// The keys are all image paths.
	allRefImagesMap, errs := helper.ScanToFixImgRelPath(docFolder, imgFolder, doFix)
	if errs != nil {
		helper.PrintAggregateError(errs)
		return
	}

	fmt.Println("Starting to scan images..")

	// Delete no reference images.
	if errs := helper.DelNoRefImags(imgFolder, allRefImagesMap, doDel); errs != nil {
		helper.PrintAggregateError(errs)
		return
	}

	fmt.Println("All done.")
}
