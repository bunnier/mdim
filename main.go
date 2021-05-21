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
		if err != helper.ExitForHelper {
			println(err.Error())
		}
		return
	}

	fmt.Println("Starting to scan markdown document..")

	// The key is all images list.
	allRefImagesMap, errs := helper.WalkDocFolderToFix(docFolder, imgFolder, doFix)
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
