package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"mdim/core/types"
)

// Iterate imageFolder to find & delete no reference images.
func DelNoRefImgs(imgFolder string, allRefImgsSet types.Set, doImgDel bool) types.AggregateError {
	var imgs []fs.DirEntry
	var err error
	if imgs, err = os.ReadDir(imgFolder); err != nil {
		err := fmt.Errorf("images: open folder failed %s %w", imgFolder, err)
		return types.NewAggregateError().AddError(err)
	}

	errCh := make(chan error) // error channel
	wg := sync.WaitGroup{}

	for _, img := range imgs {
		imgName := img.Name()

		if allRefImgsSet.Exist(imgName) {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			imgFullPath := filepath.Join(imgFolder, imgName)
			if doImgDel {
				if err := os.Remove(imgFullPath); err != nil {
					// Pass error to main goroutine.
					errCh <- fmt.Errorf("images: delete file failed %s %w", imgFullPath, err)
				} else {
					fmt.Println("images: deleted successfully", imgFullPath)
				}
			} else {
				fmt.Println("images: find a no reference image", imgFullPath)
			}
		}()
	}

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
