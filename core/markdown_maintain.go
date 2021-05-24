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

// Input=Markdown line, Group1=img title, Group2=img path, Group3=protocol, Group4=img filename
var imgTagRegexp *regexp.Regexp = regexp.MustCompile(`!\[([^]]*)]\(((?:(http[s]?://|ftp://)|[\\\./]*)(?:(?:[^\\/\n]+[\\/])*)([^\\/\n]+\.[a-zA-Z]{1,5}))\)`)

// Input=Relative path, Group1=First named folder name, Group2=relative path in imgFolder
var imgPathRegexp *regexp.Regexp = regexp.MustCompile(`^(?:(?:\.{1,2}[/\\])+)([^/\\\n]+)?[/\\](.+)$`)

// MaintainImageTags Scan docs in docFolder to fix image relative path.
// The first return is all the reference image paths Set.
func MaintainImageTags(absDocFolder string, absImgFolder string, doFix bool) (types.Set, []types.MarkdownHandleResult) {
	handleResultCh := make(chan types.MarkdownHandleResult)
	imgPathCh := make(chan types.Set)
	wg := sync.WaitGroup{}

	count := 0 // The count of handling files.
	filepath.WalkDir(absDocFolder, func(docPath string, d os.DirEntry, err error) error {
		// Just deal with markdown docs.
		if d.IsDir() || !strings.HasSuffix(docPath, ".md") {
			return nil
		}

		count++
		wg.Add(1)
		go func() {
			defer wg.Done()

			refImgsAbsPathSet, handleResult := maintainImageTagsForSingleFile(docPath, absImgFolder, doFix)
			// Pass results to main goroutine.
			handleResultCh <- handleResult
			imgPathCh <- refImgsAbsPathSet
		}()

		return nil
	})

	// Waiting for all goroutine done to close channel.
	go func() {
		wg.Wait()
		close(handleResultCh)
		close(imgPathCh)
	}()

	handleResults := make([]types.MarkdownHandleResult, 0, count)
	allRefImgsAbsPathSet := types.NewSet(count * 3)

	chOpen := true
	var handleResult types.MarkdownHandleResult
	var imgSet types.Set
	for {
		// Receive error & found image path.
		select {
		case handleResult, chOpen = <-handleResultCh:
			if chOpen {
				handleResults = append(handleResults, handleResult)
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

	return allRefImgsAbsPathSet, handleResults
}

// Fix the image urls of the doc.
// The first return is all the reference image paths Set.
func maintainImageTagsForSingleFile(docPath string, absImgFolder string, doRelPathFix bool) (types.Set, types.MarkdownHandleResult) {
	byteStream := bytes.Buffer{}          // Put the fixed text.
	refImgsAbsPathSet := types.NewSet(10) // Store all the reference image paths.
	handleResult := types.MarkdownHandleResult{DocPath: docPath}

	fileInfo, err := os.Lstat(docPath) // get perm
	if err != nil {
		handleResult.Err = fmt.Errorf("lstat file failed\n %w", err)
		return nil, handleResult
	}
	filePerm := fileInfo.Mode().Perm() // file perm

	file, err := os.OpenFile(docPath, os.O_RDWR, filePerm)
	if err != nil {
		handleResult.Err = fmt.Errorf("open file failed\n %w", err)
		return nil, handleResult
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
			handleResult.Err = fmt.Errorf("reading failed\n %w", err)
			return nil, handleResult
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

			// fix rel path
			if fixedPath, absFixedPath, err := getFixImgRelPath(docPath, imgPath, absImgFolder); err == nil {
				refImgsAbsPathSet.Add(absFixedPath)
				return fmt.Sprintf("![%s](%s)", imgTitle, fixedPath)
			} else {
				if handleResult.RelPathCannotFixedErr == nil {
					handleResult.RelPathCannotFixedErr = make([]error, 0, 1)
				}
				handleResult.RelPathCannotFixedErr = append(handleResult.RelPathCannotFixedErr, err)
				return imgTag
			}
		})

		byteStream.WriteString(fixedLine)
		handleResult.HasErrImgRelPath = handleResult.HasErrImgRelPath || fixedLine != line
	}
	file.Close()

	if !handleResult.HasErrImgRelPath || !doRelPathFix {
		return refImgsAbsPathSet, handleResult
	}

	// Write fixed content to original path.
	if err = overrideExistFile(docPath, byteStream.String(), filePerm); err != nil {
		handleResult.Err = fmt.Errorf("writing failed\n %w", err)
		return refImgsAbsPathSet, handleResult
	}

	handleResult.FixedErrImgRelPath = true
	return refImgsAbsPathSet, handleResult
}

// Try to get a correct image relative path.
//
// Handling logic:
//   1. If the imgPath is not a relative path, return error.
//   2. When the relative path is existed, return it self.
//   3. If the first named folder of imgPath do not equals the folder name of absImgFolder, return error.
//   4. Trying to point the path part of imgPath followed first named folder to the absImgFolder, than get a new path.
//   5. When the path from step4 is existed, it will be return, otherwise will return error.
//
// Return value = (relative path, abs path, error)
func getFixImgRelPath(docPath string, imgPath string, absImgFolder string) (string, string, error) {
	imgFolderName := filepath.Base(absImgFolder)
	matches := imgPathRegexp.FindAllStringSubmatch(imgPath, -1)

	// Can only handle the path that related to absImgFolder.
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

func overrideExistFile(docPath string, content string, filePerm fs.FileMode) error {
	file, err := os.OpenFile(docPath, os.O_RDWR|os.O_TRUNC, filePerm)
	if err != nil {
		return errors.New("writing open failed")
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return errors.New("writing failed")
	}
	return nil
}
