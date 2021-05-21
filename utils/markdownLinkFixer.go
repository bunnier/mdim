package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// Scan docs in docFolder to fix image relative path.
// The first return map's keys are all reference images paths
func ScanToFixImgRelPath(docFolder string, imgFolder string, doFix bool) (map[string]interface{}, []error) {
	var errCh chan error = make(chan error)
	var imagePathCh chan string = make(chan string)
	var wg sync.WaitGroup = sync.WaitGroup{}

	_ = filepath.WalkDir(docFolder, func(docPath string, d os.DirEntry, err error) error {
		// Just deal with markdown docs.
		if d.IsDir() || !strings.HasSuffix(docPath, ".md") {
			return nil
		}

		wg.Add(1)
		go func(imagePathCh chan string, errCh chan error, wg *sync.WaitGroup) {
			defer wg.Done()
			if imageSlice, err := FixImgRelPath(docPath, imgFolder, doFix); err != nil {
				errCh <- err
			} else {
				for _, v := range imageSlice {
					imagePathCh <- v
				}
			}
		}(imagePathCh, errCh, &wg)

		return nil
	})

	allRefImagesMap := make(map[string]interface{}, 100)
	aggregateErr := make([]error, 0)

	// Waiting for all goroutine done to close channel.
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(errCh)
		close(imagePathCh)
	}(&wg)

	chOpen := true
	var err error
	var imagePath string
	for {
		select {
		case err, chOpen = <-errCh:
			if chOpen {
				aggregateErr = append(aggregateErr, err)
			}
		case imagePath, chOpen = <-imagePathCh:
			if chOpen {
				allRefImagesMap[imagePath] = nil
			}
		}

		if !chOpen {
			break
		}
	}

	if len(errCh) == 0 {
		return allRefImagesMap, nil
	} else {
		return nil, aggregateErr
	}
}

// Fix the image urls of the doc.
// The first return is all the image paths slice.
func FixImgRelPath(docPath string, imageFolder string, doFix bool) ([]string, error) {
	imageTagRe := regexp.MustCompile(`!\[([^]]*)]\((?:[\\\./]*(?:(?:[^\\/\n]+[\\/])*)([^\\/\n]+\.png))\)`)

	var changed bool = false
	var byteStream bytes.Buffer            // Put the fixed text.
	imgPathsSlice := make([]string, 0, 20) // result
	var filePerm fs.FileMode
	if fileInfo, err := os.Lstat(docPath); err != nil {
		return nil, fmt.Errorf("docs: open failed %s %w", docPath, err)
	} else {
		filePerm = fileInfo.Mode().Perm()
	}

	if file, err := os.OpenFile(docPath, os.O_RDWR, filePerm); err != nil {
		return nil, fmt.Errorf("docs: open failed %s %w", docPath, err)
	} else {
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			if line, err := reader.ReadString('\n'); err != nil {
				if err == io.EOF {
					byteStream.WriteString(line)
					break
				}
				return nil, fmt.Errorf("docs: reading failed %s %w", docPath, err)
			} else {
				// Do single line replace.
				var replaceErr error
				newline := imageTagRe.ReplaceAllStringFunc(line, func(m string) string {
					matchPart := imageTagRe.FindStringSubmatch(m)
					imageTag := matchPart[0]      // whole image tag
					imageTitle := matchPart[1]    // tag title
					imageFileName := matchPart[2] // filename

					imageAbsPath := filepath.Join(imageFolder, imageFileName)
					docParentPath := filepath.Dir(docPath)

					if relPath, err := filepath.Rel(docParentPath, imageAbsPath); err != nil {
						replaceErr = fmt.Errorf("docs: calcute relative failed, from %s to %s %w", docParentPath, imageAbsPath, err)
						return m
					} else {
						imgPathsSlice = append(imgPathsSlice, imageFileName) // Add path to result.
						newLine := fmt.Sprintf("![%s](%s)", imageTitle, relPath)
						changed = changed || newLine != imageTag
						return newLine
					}
				})
				if replaceErr != nil {
					return nil, replaceErr
				}

				if _, err := byteStream.WriteString(newline); err != nil {
					_ = fmt.Errorf("docs: write fixed string error %s, %w", docPath, err)
				}
			}
		}
		file.Close()
	}

	// Write result to original path.
	if changed {
		if doFix {
			if file, err := os.OpenFile(docPath, os.O_RDWR|os.O_TRUNC, filePerm); err != nil {
				return nil, fmt.Errorf("docs: writing open failed %s %w", docPath, err)
			} else {
				defer file.Close()
				if _, err := file.WriteString(byteStream.String()); err != nil {
					return nil, fmt.Errorf("docs: writing failed %s %w", docPath, err)
				}
				fmt.Println("docs: fixed successfully", docPath)
			}
		} else {
			fmt.Println("docs: find a document with error relative path", docPath)
		}
	}

	return imgPathsSlice, nil
}
