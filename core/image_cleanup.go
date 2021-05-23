package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"mdim/core/types"
)

// Iterate imageFolder to find & delete no reference images.
func DeleteNoRefImgs(absImgFolder string, allRefImgsAbsPathSet types.Set, doImgDel bool) []types.ImageHandleResult {
	wg := sync.WaitGroup{}
	count := 0
	handleResultCh := make(chan types.ImageHandleResult)

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
					// Pass error to main goroutine.
					handleResult.Err = fmt.Errorf("images: delete no referemce image failed %s %w", imgPath, err)
					handleResult.Deleted = false
				} else {
					handleResult.Deleted = true
				}
			}
			handleResultCh <- handleResult
		}()

		return nil
	})

	// Waiting for all goroutine done to close channel.
	go func() {
		wg.Wait()
		close(handleResultCh)
	}()

	// handle result receiver
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
