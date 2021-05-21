package core

import (
	"fmt"
	"mdic/core/types"
	"os"
	"path/filepath"
	"sync"
)

// Iterate imageFolder to find & delete no reference images.
func DelNoRefImgs(imgFolder string, referenceMap map[string]interface{}, doImgDel bool) types.AggregateError {
	if imgs, err := os.ReadDir(imgFolder); err != nil {
		err := fmt.Errorf("images: open folder failed %s %w", imgFolder, err)
		return types.NewAggregateError().AddError(err)
	} else {
		errCh := make(chan error) // error channel
		wg := sync.WaitGroup{}

		for _, img := range imgs {
			wg.Add(1)
			go func(imgName string, wg *sync.WaitGroup, errCh chan error) {
				defer wg.Done()
				if _, ok := referenceMap[imgName]; !ok {
					imgFullPath := filepath.Join(imgFolder, imgName)
					if doImgDel {
						if err := os.Remove(imgFullPath); err != nil {
							errCh <- fmt.Errorf("images: delete file failed %s %w", imgFullPath, err)
						} else {
							fmt.Println("images: deleted successfully", imgFullPath)
						}
					} else {
						fmt.Println("images: find a no reference image", imgFullPath)
					}
				}
			}(img.Name(), &wg, errCh)
		}

		// Waiting for all goroutine done to close channel.
		go func(wg *sync.WaitGroup) {
			wg.Wait()
			close(errCh)
		}(&wg)

		// channel receiver
		aggErr := types.NewAggregateError()
		for {
			if err, ok := <-errCh; ok {
				aggErr.AddError(err)
			} else {
				break
			}
		}

		if aggErr.Len() == 0 {
			return nil
		} else {
			return aggErr
		}
	}
}
