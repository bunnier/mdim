package cleaner

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bunnier/mdim/internal/base"
)

// DeleteNoRefImgs iterate imageFolder to find & delete no reference images.
func DeleteNoRefImgs(absImgFolder string, allRefImgsAbsPathSet base.Set) []HandleResult {
	handleResultCh := make(chan HandleResult)
	count := 0 // The count of handling files.
	filepath.WalkDir(absImgFolder, func(imgPath string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if allRefImgsAbsPathSet.Exist(imgPath) {
			return nil
		}

		count++
		go func() {
			handleResult := HandleResult{ImagePath: imgPath}
			defer func() {
				// Ensure to pass handling result to main goroutine, otherwise can cause deadlock.
				handleResultCh <- handleResult
			}()

			if err := os.Remove(imgPath); err != nil {
				handleResult.Err = fmt.Errorf("delete no referemce image failed\n%w", err)
			}
		}()

		return nil
	})

	// handling result receiver
	handleResultSlice := make([]HandleResult, 0, count)
	for count > 0 {
		handleResult := <-handleResultCh
		handleResultSlice = append(handleResultSlice, handleResult)
		count--
	}
	close(handleResultCh)

	return handleResultSlice
}
