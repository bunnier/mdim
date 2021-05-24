package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"mdim/core/types"
)

// DeleteNoRefImgs Iterate imageFolder to find & delete no reference images.
func DeleteNoRefImgs(absImgFolder string, allRefImgsAbsPathSet types.Set, doImgDel bool) []types.ImageHandleResult {
	handleResultCh := make(chan types.ImageHandleResult)
	wg := sync.WaitGroup{}
	count := 0 // The count of handling files.
	filepath.WalkDir(absImgFolder, func(imgPath string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if allRefImgsAbsPathSet.Exist(imgPath) {
			return nil
		}

		wg.Add(1)
		count++
		go func() {
			defer wg.Done()
			handleResult := types.ImageHandleResult{ImagePath: imgPath}
			if doImgDel {
				if err := os.Remove(imgPath); err != nil {
					handleResult.Err = fmt.Errorf("delete no referemce image failed\n%w", err)
					handleResult.Deleted = false
				} else {
					handleResult.Deleted = true
				}
			}
			// Pass handling result to main goroutine.
			handleResultCh <- handleResult
		}()

		return nil
	})

	// Waiting for all goroutine done to close channel.
	go func() {
		wg.Wait()
		close(handleResultCh)
	}()

	// handling result receiver
	handleResultSlice := make([]types.ImageHandleResult, 0, count)
	for {
		handleResult, chOpen := <-handleResultCh
		if !chOpen {
			break
		}
		handleResultSlice = append(handleResultSlice, handleResult)
	}

	return handleResultSlice
}
