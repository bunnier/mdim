package core

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

	"mdim/core/types"
)

// Scan docs in docFolder to fix image relative path.
// The first return map's keys are all reference images paths.
func ScanMarkdownFiles(docFolder string, imgFolder string, doFix bool) (types.Set, types.AggregateError) {
	errCh := make(chan error)
	imgPathCh := make(chan string)
	wg := sync.WaitGroup{}

	// Error will pass to errCh.
	filepath.WalkDir(docFolder, func(docPath string, d os.DirEntry, err error) error {
		// Just deal with markdown docs.
		if d.IsDir() || !strings.HasSuffix(docPath, ".md") {
			return nil
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if imgPathSlice, err := scanMarkdownFile(docPath, imgFolder, doFix); err != nil {
				errCh <- err // Pass error to main goroutine.
			} else {
				for _, v := range imgPathSlice {
					imgPathCh <- v // Pass found image paths to main goroutine.
				}
			}
		}()

		return nil
	})

	allRefImgsSet := types.NewSet(100)
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
				allRefImgsSet.Put(imgPath)
			}
		}

		if !chOpen {
			break
		}
	}

	if aggErr.Len() == 0 {
		return allRefImgsSet, nil
	} else {
		return nil, aggErr
	}
}

// Fix the image urls of the doc.
// The first return is all the image paths slice.
func scanMarkdownFile(docPath string, imgFolder string, doRelPathFix bool) ([]string, error) {
	changed := false
	byteStream := bytes.Buffer{}           // Put the fixed text.
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

		// Search image tags.
		refPathSet := findImgInLine(line)

		// image not found OR do not fix
		if refPathSet.IsEmpty() || !doRelPathFix {
			byteStream.WriteString(line)
			continue
		}

		// Do fix.
		if fixedLine, err := fixSingleLineImgRelPat(line, docPath, imgFolder); err != nil {
			return nil, fmt.Errorf("docs: image relative path fixing error %s %w", docPath, err)
		} else {
			byteStream.WriteString(fixedLine)
			if fixedLine != line {
				changed = changed || true
			}
		}
	}
	file.Close()

	if !changed || !doRelPathFix {
		return imgPathsSlice, nil
	}

	// Write fixed content to original path.
	if err = writeFixedContent(docPath, byteStream.String(), filePerm); err != nil {
		return nil, err
	}

	fmt.Println("docs: fixed successfully", docPath)
	return imgPathsSlice, nil
}

// Group1=img title, Group2=img path, Group3=protocol, Group4=img filename
var imgTagRegexp *regexp.Regexp = regexp.MustCompile(`!\[([^]]*)]\(((?:(http[s]?://|ftp://)|[\\\./]*)(?:(?:[^\\/\n]+[\\/])*)([^\\/\n]+\.[a-zA-Z]{1,5}))\)`)

// Scan single line to find imgpath.
func findImgInLine(line string) types.Set {
	refImgsSet := types.NewSet(2)

	for _, matchParts := range imgTagRegexp.FindAllStringSubmatch(line, -1) {
		refImgsSet.Add(matchParts[2]) // 2 is img path.
	}

	return refImgsSet
}

func fixSingleLineImgRelPat(line string, docPath string, imgFolder string) (string, error) {
	// Do fix.
	var replaceErr error
	fixedLine := imgTagRegexp.ReplaceAllStringFunc(line, func(imgTag string) string {
		matchParts := imgTagRegexp.FindStringSubmatch(imgTag) // matchLine is whole image tag
		imgTitle := matchParts[1]                             // tag title
		imgProtocol := matchParts[3]                          // protocol
		imgFileName := matchParts[4]                          // filename

		// Do not deal with remote imgs here.
		if imgProtocol != "" {
			return imgTag
		}

		docFolder := filepath.Dir(docPath) // folder of current doc
		imgAbsPath := filepath.Join(imgFolder, imgFileName)

		relPath, err := filepath.Rel(docFolder, imgAbsPath)
		if err != nil {
			replaceErr = fmt.Errorf("docs: calcute relative failed, from %s to %s %w", docFolder, imgAbsPath, err)
			return imgTag
		}

		newImgTag := fmt.Sprintf("![%s](%s)", imgTitle, relPath)
		return newImgTag
	})

	if replaceErr != nil {
		return "", replaceErr
	}

	return fixedLine, nil
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
