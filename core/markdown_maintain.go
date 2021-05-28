package core

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"mdim/core/types"
)

var (
	// Input=Markdown line, Group1=img title, Group2=img path, Group3=protocol, Group4=img filename
	imgTagRegexp *regexp.Regexp = regexp.MustCompile(`!\[([^]]*)]\(((?:(http[s]?://|ftp://)|[\\\./]*)(?:(?:[^\\/\n]+[\\/])*)([^\\/\n]+\.[a-zA-Z]{1,5}))\)`)

	// Input=Relative path, Group1=First named folder name, Group2=relative path in imgFolder
	imgPathRegexp *regexp.Regexp = regexp.MustCompile(`^(?:(?:\.{1,2}[/\\])+)([^/\\\n]+)?[/\\](.+)$`)

	// For slash replace.
	slashReplaceRegexp = regexp.MustCompile(`\\`)

	// For valid http/https.
	httpValidRegex = regexp.MustCompile(`^http[s]://`)

	// For get web image suffix.
	httpImgRegex = regexp.MustCompile(`^(?:http[s]?://)(?:[^/]+/)+.+(\.[a-zA-Z]{1,5})$`)
)

// MaintainImageTags Scan docs in docFolder to fix image relative path.
// The first return is all the reference image paths Set.
func MaintainImageTags(absDocFolder string, absImgFolder string, doSave bool, doWebImgDownload bool) (types.Set, []types.MarkdownHandleResult) {
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

			refImgsAbsPathSet, handleResult := maintainImageTagsForSingleFile(docPath, absImgFolder, doSave, doWebImgDownload)
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
func maintainImageTagsForSingleFile(docPath string, absImgFolder string, doSave bool, doWebImgDownload bool) (types.Set, types.MarkdownHandleResult) {
	refImgsAbsPathSet := types.NewSet(10) // Store all the reference image paths.
	handleResult := types.MarkdownHandleResult{DocPath: docPath}

	// get doc file content
	contentBytes, err := os.ReadFile(docPath)
	if err != nil {
		handleResult.Err = fmt.Errorf("reading failed\n%w", err)
		return nil, handleResult
	}
	content := string(contentBytes)

	// line workflow
	fixedContent := imgTagRegexp.ReplaceAllStringFunc(content, func(imgTag string) string {
		matchParts := imgTagRegexp.FindStringSubmatch(imgTag) // matchLine is whole image tag
		imgTitle := matchParts[1]                             // tag title
		imgPath := matchParts[2]                              // img path
		imgProtocol := strings.ToLower(matchParts[3])         // protocol

		// Can't handle this url.
		if imgProtocol != "" && (!doWebImgDownload || !httpValidRegex.MatchString(imgProtocol)) {
			refImgsAbsPathSet.Add(imgPath)
			return imgTag
		}

		// Follow logics can be convert to a pipe~

		if imgProtocol != "" {
			// handle web url
			if fixedPath, absFixedPath, err := convertRemoteImageToLocal(docPath, imgPath, absImgFolder); err == nil {
				refImgsAbsPathSet.Add(absFixedPath)
				return fmt.Sprintf("![%s](%s)", imgTitle, fixedPath)
			} else {
				if handleResult.WebImgDownloadErr == nil {
					handleResult.WebImgDownloadErr = make([]error, 0, 1)
				}
				handleResult.WebImgDownloadErr = append(handleResult.WebImgDownloadErr, err)
				return imgTag
			}
		} else {
			// handle local url
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
		}
	})

	handleResult.HasChangeDuringMaintain = fixedContent != content

	if !handleResult.HasChangeDuringMaintain || !doSave {
		return refImgsAbsPathSet, handleResult
	}

	// Write fixed content to original path.
	if err = overrideExistFile(docPath, fixedContent); err != nil {
		handleResult.Err = fmt.Errorf("writing failed\n%w", err)
		return refImgsAbsPathSet, handleResult
	}

	handleResult.SavedMaintainResult = true
	return refImgsAbsPathSet, handleResult
}

// Download the web images then return the relative path to docPath and absPath
// Return value = (relative path, abs path, error)
func convertRemoteImageToLocal(docPath string, imgPath string, absImgFolder string) (string, string, error) {
	matches := httpImgRegex.FindAllStringSubmatch(imgPath, -1)
	// Can only handle the path that related to absImgFolder.
	if len(matches) != 1 || len(matches[0]) != 2 {
		return "", "", errors.New("can not handle this url")
	}

	imgSuffix := matches[0][1]

	fmt.Println("Begin to download web img:", imgPath)
	httpClient := http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := httpClient.Get(imgPath)
	fmt.Println("Img downloaded:", imgPath)

	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("status code=%s, url=%s", resp.Status, imgPath)
	}

	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	filename := time.Now().Format("2006-01-2-15-04-05.000") + imgSuffix
	absImgPath := filepath.Join(absImgFolder, filename)
	if err := os.WriteFile(absImgPath, imgBytes, 666); err != nil {
		return "", "", err
	}

	currentDocFolder := filepath.Dir(docPath)
	fixImgRelPath, err := filepath.Rel(currentDocFolder, absImgPath)
	if err != nil {
		return "", "", err
	}

	return fixImgRelPath, absImgPath, nil
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

	// Unify to forward slash.
	fixImgRelPath = slashReplaceRegexp.ReplaceAllString(fixImgRelPath, "/")

	return fixImgRelPath, fixImgAbsPath, nil
}

func overrideExistFile(docPath string, content string) error {
	fileInfo, err := os.Lstat(docPath) // get perm
	if err != nil {
		return fmt.Errorf("lstat file failed\n%w", err)
	}
	filePerm := fileInfo.Mode().Perm() // file perm

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
