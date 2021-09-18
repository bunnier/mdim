package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"unsafe"

	"github.com/bunnier/mdim/internal/types"
)

var (
	// Input=Markdown line, Group1=img title, Group2=img path, Group3=protocol, Group4=img filename
	imgTagRegexp *regexp.Regexp = regexp.MustCompile(`!\[([^]]*)]\(((?:(http[s]?://|ftp://)|[\\\./]*)(?:(?:[^\\/\n]+[\\/])*)([^\\/\n]+\.[a-zA-Z]{1,5}))\)`)

	// For valid http/https.
	httpValidRegex = regexp.MustCompile(`^http[s]://`)
)

// MarkdownImageTag represent an image tag in markdown document.
type MarkdownImageTag struct {
	IsWebUrl     bool
	Tag          string
	Title        string
	DocPath      string
	ImgPath      string
	AbsImgFolder string
	Protocal     string
}

// ImageMaintainStep is markdown handle step.
type ImageMaintainStep func(imgTag *MarkdownImageTag, handleResult types.MarkdownHandleResult) error

// WalkDirToHandleDocs will scan docs in docFolder to fix image relative path.
// The first return is all the reference image paths Set.
func WalkDirToHandleDocs(absDocFolder string, absImgFolder string, doSave bool, doWebImgDownload bool) []types.MarkdownHandleResult {
	handleResultCh := make(chan types.MarkdownHandleResult)
	wg := sync.WaitGroup{}

	fileNum := 0 // The count of handling files.
	filepath.WalkDir(absDocFolder, func(docPath string, d os.DirEntry, err error) error {
		// Just deal with markdown docs.
		if d.IsDir() || !strings.HasSuffix(docPath, ".md") {
			return nil
		}

		fileNum++
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleResultCh <- handleDoc(docPath, absImgFolder, doSave, doWebImgDownload, TryDownloadImage, TryFixLocalImageRelpath)
		}()

		return nil
	})

	// Waiting for all goroutine done to close channel.
	go func() {
		wg.Wait()
		close(handleResultCh)
	}()

	aggreagateResult := make([]types.MarkdownHandleResult, 0, fileNum)
	for {
		handleResult, chOpen := <-handleResultCh
		if !chOpen {
			break
		}
		aggreagateResult = append(aggreagateResult, handleResult)
	}

	return aggreagateResult
}

// Fix the image urls of the doc.
// The first return is all the reference image paths Set.
func handleDoc(docPath string, absImgFolder string, doSave bool, doWebImgDownload bool, steps ...ImageMaintainStep) types.MarkdownHandleResult {
	handleResult := types.MarkdownHandleResult{DocPath: docPath}
	// get doc file content
	contentBytes, err := os.ReadFile(docPath)
	if err != nil {
		handleResult.Err = fmt.Errorf("reading failed\n%w", err)
		return handleResult
	}

	handleResult.AllRefImgs = types.NewSet(10) // To store reference image paths.

	// directly convert for saving memory
	content := *(*string)(unsafe.Pointer(&contentBytes))

	// line workflow
	fixedContent := imgTagRegexp.ReplaceAllStringFunc(content, func(imgTag string) string {
		matchParts := imgTagRegexp.FindStringSubmatch(imgTag) // matchLine is whole image tag
		imgTitle := matchParts[1]                             // tag title
		imgPath := matchParts[2]                              // img path
		imgProtocol := strings.ToLower(matchParts[3])         // protocol

		// Can't handle this url.
		if imgProtocol != "" && (!doWebImgDownload || !httpValidRegex.MatchString(imgProtocol)) {
			handleResult.AllRefImgs.Add(imgPath)
			return imgTag
		}

		imgTagInfo := &MarkdownImageTag{
			IsWebUrl:     imgProtocol != "",
			Tag:          imgTag,
			Title:        imgTitle,
			DocPath:      docPath,
			ImgPath:      imgPath,
			AbsImgFolder: absImgFolder,
			Protocal:     imgProtocol,
		}

		for _, handleStep := range steps {
			handleStep(imgTagInfo, handleResult)
		}

		return imgTagInfo.Tag
	})

	handleResult.HasChangeDuringMaintain = fixedContent != content

	if !handleResult.HasChangeDuringMaintain || !doSave {
		return handleResult
	}

	// Write fixed content to original path.
	if err = overrideExistFile(docPath, fixedContent); err != nil {
		handleResult.Err = fmt.Errorf("writing failed\n%w", err)
		return handleResult
	}

	handleResult.SavedMaintainResult = true
	return handleResult
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
