package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"mdim/core/types"
)

// Iterate imageFolder to find & delete no reference images.
func DelNoRefImgs(absImgFolder string, allRefImgsAbsPathSet types.Set, doImgDel bool) types.AggregateError {

	errCh := make(chan error) // error channel
	wg := sync.WaitGroup{}

	filepath.WalkDir(absImgFolder, func(imgPath string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if allRefImgsAbsPathSet.Exist(imgPath) {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if doImgDel {
				if err := os.Remove(imgPath); err != nil {
					// Pass error to main goroutine.
					errCh <- fmt.Errorf("images: delete no referemce image failed %s %w", imgPath, err)
				} else {
					fmt.Println("images: delete a no reference image successfully", imgPath)
				}
			} else {
				fmt.Println("images: find a no reference image", imgPath)
			}
		}()

		return nil
	})

	// Waiting for all goroutine done to close channel.
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// error receiver
	aggErr := types.NewAggregateError()
	for {
		err, chOpen := <-errCh
		if !chOpen {
			break
		}
		aggErr.AddError(err)
	}

	if aggErr.Len() == 0 {
		return nil
	} else {
		return aggErr
	}
}
