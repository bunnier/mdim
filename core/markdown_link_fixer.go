package core

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"mdim/core/types"
)

// Scan docs in docFolder to fix image relative path.
// The first return map's keys are all reference images paths.
func ScanMarkdownFiles(absDocFolder string, absImgFolder string, doFix bool) (types.Set, types.AggregateError) {
	errCh := make(chan error)
	imgPathCh := make(chan types.Set)
	wg := sync.WaitGroup{}

	// Error will pass to errCh.
	filepath.WalkDir(absDocFolder, func(docPath string, d os.DirEntry, err error) error {
		// Just deal with markdown docs.
		if d.IsDir() || !strings.HasSuffix(docPath, ".md") {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if refImgsAbsPathSet, err := scanMarkdownFile(docPath, absImgFolder, doFix); err != nil {
				errCh <- err // Pass error to main goroutine.
			} else {
				imgPathCh <- refImgsAbsPathSet // Pass found image paths to main goroutine.
			}
		}()

		return nil
	})

	allRefImgsAbsPathSet := types.NewSet(100)
	aggErr := types.NewAggregateError()

	// Waiting for all goroutine done to close channel.
	go func() {
		wg.Wait()
		close(errCh)
		close(imgPathCh)
	}()

	chOpen := true
	var err error
	var imgSet types.Set
	for {
		// Receive error & found image path.
		select {
		case err, chOpen = <-errCh:
			if chOpen {
				aggErr.AddError(err)
			}
		case imgSet, chOpen = <-imgPathCh:
			if chOpen {
				allRefImgsAbsPathSet.Merge(imgSet)
			}
		}

		if !chOpen {
			break
		}
	}

	if aggErr.Len() == 0 {
		return allRefImgsAbsPathSet, nil
	} else {
		return nil, aggErr
	}
}

// Group1=img title, Group2=img path, Group3=protocol, Group4=img filename
var imgTagRegexp *regexp.Regexp = regexp.MustCompile(`!\[([^]]*)]\(((?:(http[s]?://|ftp://)|[\\\./]*)(?:(?:[^\\/\n]+[\\/])*)([^\\/\n]+\.[a-zA-Z]{1,5}))\)`)

// Fix the image urls of the doc.
// The first return is all the reference image paths set.
func scanMarkdownFile(docPath string, absImgFolder string, doRelPathFix bool) (types.Set, error) {
	byteStream := bytes.Buffer{}          // Put the fixed text.
	refImgsAbsPathSet := types.NewSet(10) // Store all the reference image paths.

	fileInfo, err := os.Lstat(docPath) // get perm
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
	changed := false
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				byteStream.WriteString(line)
				break
			}
			return nil, fmt.Errorf("docs: reading failed %s %w", docPath, err)
		}

		// line workflow
		fixedLine := imgTagRegexp.ReplaceAllStringFunc(line, func(imgTag string) string {
			matchParts := imgTagRegexp.FindStringSubmatch(imgTag) // matchLine is whole image tag
			imgTitle := matchParts[1]                             // tag title
			imgPath := matchParts[2]                              // img path
			imgProtocol := matchParts[3]                          // protocol

			if imgProtocol != "" { // do not handle remote image
				refImgsAbsPathSet.Add(imgPath)
				return imgTag
			}

			if fixedPath, absFixedPath, err := getFixImgRelPath(docPath, imgPath, absImgFolder); err == nil {
				// save abs path
				refImgsAbsPathSet.Add(absFixedPath)
				return fmt.Sprintf("![%s](%s)", imgTitle, fixedPath)
			} else {
				// log then continuess
				fmt.Printf("\ndocs: failed to fix this image path\n%s\n%s\n%s\n\n", docPath, imgPath, err.Error())
				return imgTag
			}
		})

		byteStream.WriteString(fixedLine)
		changed = changed || fixedLine != line
	}
	file.Close()

	if !changed || !doRelPathFix {
		return refImgsAbsPathSet, nil
	}

	// Write fixed content to original path.
	if err = writeFixedContent(docPath, byteStream.String(), filePerm); err != nil {
		return nil, err
	}

	fmt.Println("docs: fixed successfully", docPath)
	return refImgsAbsPathSet, nil
}

// Group1=imgFolder Name, Group2=relative path in imgFolder
var imgPathRegexp *regexp.Regexp = regexp.MustCompile(`^(?:(?:\.{1,2}[/\\])+)([^/\\\n]+)?[/\\](.+)$`)

// Try to fix image relative path to imgFolder
// return (relative path, abs path, error)
func getFixImgRelPath(docPath string, imgPath string, absImgFolder string) (string, string, error) {
	imgFolderName := filepath.Base(absImgFolder)
	matches := imgPathRegexp.FindAllStringSubmatch(imgPath, -1)

	// can not handle this path
	if len(matches) != 1 || len(matches[0]) != 3 || matches[0][1] != imgFolderName {
		return "", "", errors.New("can not handle this path")
	}
	relPathInImgFolder := matches[0][2]

	fixImgAbsPath := filepath.Join(absImgFolder, relPathInImgFolder)
	_, err := os.Stat(fixImgAbsPath)
	if err != nil { // can not handle
		return "", "", err
	}

	currentDocFolder := filepath.Dir(docPath)
	fixImgRelPath, err := filepath.Rel(currentDocFolder, fixImgAbsPath)
	if err != nil {
		return "", "", err
	}

	return fixImgRelPath, fixImgAbsPath, nil
}

func writeFixedContent(docPath string, content string, filePerm fs.FileMode) error {
	file, err := os.OpenFile(docPath, os.O_RDWR|os.O_TRUNC, filePerm)
	if err != nil {
		return fmt.Errorf("docs: writing open failed %s %w", docPath, err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("docs: writing failed %s %w", docPath, err)
	}
	return nil
}
