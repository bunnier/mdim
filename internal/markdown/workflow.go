package markdown

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"unsafe"

	"github.com/bunnier/mdim/internal/base"
)

var (
	// Input=Markdown line, Group1=img title, Group2=img path, Group3=protocol, Group4=img filename
	imgTagRegexp *regexp.Regexp = regexp.MustCompile(`!\[([^]]*)]\(((?:(http[s]?://|ftp://)|[\\\./]*)(?:(?:[^\\/\n]+[\\/])*)([^\\/\n]+\.[a-zA-Z]{1,5}))\)`)

	// Input=Relative path, Group1=First named folder name, Group2=relative path in imgFolder
	imgPathRegexp *regexp.Regexp = regexp.MustCompile(`^(?:(?:\.{1,2}[/\\])+)([^/\\\n]+)?[/\\](.+)$`)
)

// ImageTag represent an image tag in markdown document.
type ImageTag struct {
	IsWebUrl     bool
	WholeTag     string
	ImgTitle     string
	DocPath      string
	ImgPath      string
	AbsImgFolder string
	Protocal     string
}

// ImageMaintainStep is markdown handle step.
type ImageMaintainStep func(imgTag *ImageTag, handleResult *HandleResult) error

// WalkDirToHandleDocs will scan docs in docFolder to fix image relative path.
// The first return is all the reference image paths Set.
func WalkDirToHandleDocs(absDocPath string, absDocFolder string, absImgFolder string, steps []ImageMaintainStep) []HandleResult {
	handleResultCh := make(chan HandleResult)

	fileNum := 0          // The count of handling files.
	if absDocPath != "" { // single file workflow
		fileNum++
		go func() {
			handleResultCh <- handleDoc(absDocPath, absImgFolder, steps)
			close(handleResultCh)
		}()
	} else { // folder
		wg := sync.WaitGroup{}
		filepath.WalkDir(absDocFolder, func(docPath string, d os.DirEntry, err error) error {
			// Just deal with .md docs.
			if d.IsDir() || !strings.HasSuffix(docPath, ".md") {
				return nil
			}

			fileNum++
			wg.Add(1)
			go func() {
				defer wg.Done()
				handleResultCh <- handleDoc(docPath, absImgFolder, steps)
			}()

			return nil
		})

		// Waiting for all goroutine done to close channel.
		go func() {
			wg.Wait()
			close(handleResultCh)
		}()
	}

	aggreagateResult := make([]HandleResult, 0, fileNum)
	for handleResult := range handleResultCh {
		aggreagateResult = append(aggreagateResult, handleResult)
	}

	return aggreagateResult
}

// Fix the image urls of the doc.
// The first return is all the reference image paths Set.
func handleDoc(docPath string, absImgFolder string, steps []ImageMaintainStep) HandleResult {
	handleResult := HandleResult{DocPath: docPath}
	// get doc file content
	contentBytes, err := os.ReadFile(docPath)
	if err != nil {
		handleResult.Err = fmt.Errorf("reading failed\n%w", err)
		return handleResult
	}

	handleResult.AllRefImgs = base.NewSet(10) // To store reference image paths.

	// directly convert for saving memory
	content := *(*string)(unsafe.Pointer(&contentBytes))

	// line workflow
	fixedContent := imgTagRegexp.ReplaceAllStringFunc(content, func(wholeImgTag string) string {
		matchParts := imgTagRegexp.FindStringSubmatch(wholeImgTag) // matchLine is whole image tag
		imgTitle := matchParts[1]                                  // tag title
		imgPath := matchParts[2]                                   // img path
		imgProtocol := strings.ToLower(matchParts[3])              // protocol

		imgTagInfo := &ImageTag{
			IsWebUrl:     imgProtocol != "",
			Protocal:     imgProtocol,
			WholeTag:     wholeImgTag,
			ImgTitle:     imgTitle,
			DocPath:      docPath,
			ImgPath:      imgPath,
			AbsImgFolder: absImgFolder,
		}

		for _, handleStep := range steps {
			handleStep(imgTagInfo, &handleResult)
		}

		return imgTagInfo.WholeTag
	})

	handleResult.HasChangeDuringWorkflow = fixedContent != content

	if !handleResult.HasChangeDuringWorkflow {
		return handleResult
	}

	// Write fixed content to original path.
	if err = overrideExistFile(docPath, fixedContent); err != nil {
		handleResult.Err = fmt.Errorf("writing failed\n%w", err)
		return handleResult
	}

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
