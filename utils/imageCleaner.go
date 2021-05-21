package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Iterate imageFolder to find & delete no reference images.
func DelNoRefImags(imageFolder string, referenceMap map[string]interface{}, doDel bool) []error {
	if images, err := os.ReadDir(imageFolder); err != nil {
		return []error{fmt.Errorf("images: open folder failed %s %w", imageFolder, err)}
	} else {
		errCh := make(chan error) // error channel
		wg := sync.WaitGroup{}

		for _, image := range images {
			wg.Add(1)
			go func(imageName string, wg *sync.WaitGroup, errCh chan error) {
				defer wg.Done()
				if _, ok := referenceMap[imageName]; !ok {
					imageFullPath := filepath.Join(imageFolder, imageName)
					if doDel {
						if err := os.Remove(imageFullPath); err != nil {
							errCh <- fmt.Errorf("images: delete file failed %s %w", imageFullPath, err)
						} else {
							fmt.Println("images: deleted successfully", imageFullPath)
						}
					} else {
						fmt.Println("images: find a no reference image", imageFullPath)
					}
				}
			}(image.Name(), &wg, errCh)
		}

		// Waiting for all goroutine done to close channel.
		go func(wg *sync.WaitGroup) {
			wg.Wait()
			close(errCh)
		}(&wg)

		// channel receiver
		aggregateErr := make([]error, 0, 0) // Expect 0 error, hah~
		for {
			if err, ok := <-errCh; ok {
				aggregateErr = append(aggregateErr, err)
			} else {
				break
			}
		}

		if len(aggregateErr) == 0 {
			return nil
		} else {
			return aggregateErr
		}
	}
}
