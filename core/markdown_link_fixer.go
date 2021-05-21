package core

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"mdic/core/types"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Scan docs in docFolder to fix image relative path.
// The first return map's keys are all reference images paths.
func ScanToFixImgRelPath(docFolder string, imgFolder string, doFix bool) (map[string]interface{}, types.AggregateError) {
	var errCh chan error = make(chan error)
	var imgPathCh chan string = make(chan string)
	var wg sync.WaitGroup = sync.WaitGroup{}

	// Error will pass to errCh.
	_ = filepath.WalkDir(docFolder, func(docPath string, d os.DirEntry, err error) error {
		// Just deal with markdown docs.
		if d.IsDir() || !strings.HasSuffix(docPath, ".md") {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if imgPathSlice, err := FixImgRelPath(docPath, imgFolder, doFix); err != nil {
				errCh <- err // Pass error to main goroutine.
			} else {
				for _, v := range imgPathSlice {
					imgPathCh <- v // Pass found image paths to main goroutine.
				}
			}
		}()

		return nil
	})

	allRefImgsMap := make(map[string]interface{}, 100)
	aggErr := types.NewAggregateError()

	// Waiting for all goroutine done to close channel.
	go func() {
		wg.Wait()
		close(errCh)
		close(imgPathCh)
	}()

	chOpen := true
	var err error
	var imgPath string
	for {
		// Receive error & found image path.
		select {
		case err, chOpen = <-errCh:
			if chOpen {
				aggErr.AddError(err)
			}
		case imgPath, chOpen = <-imgPathCh:
			if chOpen {
				allRefImgsMap[imgPath] = nil
			}
		}

		if !chOpen {
			break
		}
	}

	if aggErr.Len() == 0 {
		return allRefImgsMap, nil
	} else {
		return nil, aggErr
	}
}

// Fix the image urls of the doc.
// The first return is all the image paths slice.
func FixImgRelPath(docPath string, imgFolder string, doRelPathFix bool) ([]string, error) {
	imgTagRe := regexp.MustCompile(`!\[([^]]*)]\((?:[\\\./]*(?:(?:[^\\/\n]+[\\/])*)([^\\/\n]+\.png))\)`)

	var changed bool = false
	var byteStream bytes.Buffer            // Put the fixed text.
	imgPathsSlice := make([]string, 0, 20) // result

	fileInfo, err := os.Lstat(docPath) // to get perm
	if err != nil {
		return nil, fmt.Errorf("docs: open failed %s %w", docPath, err)
	}
	filePerm := fileInfo.Mode().Perm() // file perm

	file, err := os.OpenFile(docPath, os.O_RDWR, filePerm)
	if err != nil {
		return nil, fmt.Errorf("docs: open failed %s %w", docPath, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				byteStream.WriteString(line)
				break
			}
			return nil, fmt.Errorf("docs: reading failed %s %w", docPath, err)
		}

		// Do single line replace.
		var replaceErr error
		newline := imgTagRe.ReplaceAllStringFunc(line, func(m string) string {
			matchPart := imgTagRe.FindStringSubmatch(m)
			imgTag := matchPart[0]      // whole image tag
			imgTitle := matchPart[1]    // tag title
			imgFileName := matchPart[2] // filename

			imgAbsPath := filepath.Join(imgFolder, imgFileName)
			docParentPath := filepath.Dir(docPath)

			relPath, err := filepath.Rel(docParentPath, imgAbsPath)
			if err != nil {
				replaceErr = fmt.Errorf("docs: calcute relative failed, from %s to %s %w", docParentPath, imgAbsPath, err)
				return m
			}
			imgPathsSlice = append(imgPathsSlice, imgFileName) // Add path to result.
			newLine := fmt.Sprintf("![%s](%s)", imgTitle, relPath)
			changed = changed || newLine != imgTag
			return newLine
		})

		if replaceErr != nil {
			return nil, replaceErr
		}

		if _, err := byteStream.WriteString(newline); err != nil {
			return nil, fmt.Errorf("docs: write fixed string error %s, %w", docPath, err)
		}
	}
	file.Close()

	if !changed || !doRelPathFix {
		return imgPathsSlice, nil
	}

	// Write result to original path.
	file, err = os.OpenFile(docPath, os.O_RDWR|os.O_TRUNC, filePerm)
	if err != nil {
		return nil, fmt.Errorf("docs: writing open failed %s %w", docPath, err)
	}
	defer file.Close()

	if _, err := file.WriteString(byteStream.String()); err != nil {
		return nil, fmt.Errorf("docs: writing failed %s %w", docPath, err)
	}
	fmt.Println("docs: fixed successfully", docPath)
	return imgPathsSlice, nil
}
