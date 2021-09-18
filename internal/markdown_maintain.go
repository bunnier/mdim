package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/bunnier/mdim/internal/types"
)

var (
	// Input=Relative path, Group1=First named folder name, Group2=relative path in imgFolder
	imgPathRegexp *regexp.Regexp = regexp.MustCompile(`^(?:(?:\.{1,2}[/\\])+)([^/\\\n]+)?[/\\](.+)$`)

	// For get web image suffix.
	httpImgRegex = regexp.MustCompile(`^(?:http[s]?://)(?:[^/]+/)+.+(\.[a-zA-Z]{1,5})$`)

	// For slash replace.
	slashReplaceRegexp = regexp.MustCompile(`\\`)
)

// MaintainImageRelateUrl
func MaintainImageRelateUrl(imgTag *MarkdownImageTag, handleResult types.MarkdownHandleResult) error {
	if fixedPath, absFixedPath, err := getFixImgRelpath(imgTag.DocPath, imgTag.ImgPath, imgTag.AbsImgFolder); err != nil {
		if handleResult.RelPathCannotFixedErr == nil {
			handleResult.RelPathCannotFixedErr = make([]error, 0, 1)
		}
		handleResult.RelPathCannotFixedErr = append(handleResult.RelPathCannotFixedErr, err)
		return err
	} else {
		handleResult.AllRefImgs.Add(absFixedPath)
		imgTag.Tag = fmt.Sprintf("![%s](%s)", imgTag.Title, fixedPath)
		return nil
	}
}

// TryDownloadImage
func TryDownloadImage(imgTag *MarkdownImageTag, handleResult types.MarkdownHandleResult) error {
	if !imgTag.IsWebUrl {
		return nil // No suit.
	}

	if fixedPath, absFixedPath, err := convertRemoteImageToLocal(imgTag.DocPath, imgTag.ImgPath, imgTag.AbsImgFolder); err != nil {
		if handleResult.WebImgDownloadErr == nil {
			handleResult.WebImgDownloadErr = make([]error, 0, 1)
		}
		handleResult.WebImgDownloadErr = append(handleResult.WebImgDownloadErr, err)
		return err
	} else {
		handleResult.AllRefImgs.Add(absFixedPath)
		imgTag.Tag = fmt.Sprintf("![%s](%s)", imgTag.Title, fixedPath)
		return nil
	}
}

// TryFixLocalImageRelpath
func TryFixLocalImageRelpath(imgTag *MarkdownImageTag, handleResult types.MarkdownHandleResult) error {
	if imgTag.IsWebUrl {
		return nil // No suit.
	}

	if fixedPath, absFixedPath, err := getFixImgRelpath(imgTag.DocPath, imgTag.ImgPath, imgTag.AbsImgFolder); err != nil {
		if handleResult.RelPathCannotFixedErr == nil {
			handleResult.RelPathCannotFixedErr = make([]error, 0, 1)
		}
		handleResult.RelPathCannotFixedErr = append(handleResult.RelPathCannotFixedErr, err)
		return err
	} else {
		handleResult.AllRefImgs.Add(absFixedPath)
		imgTag.Tag = fmt.Sprintf("![%s](%s)", imgTag.Title, fixedPath)
		return nil
	}
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
func getFixImgRelpath(docPath string, imgPath string, absImgFolder string) (string, string, error) {
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
	if err := os.WriteFile(absImgPath, imgBytes, 0666); err != nil {
		return "", "", err
	}

	currentDocFolder := filepath.Dir(docPath)
	fixImgRelPath, err := filepath.Rel(currentDocFolder, absImgPath)
	if err != nil {
		return "", "", err
	}

	fixImgRelPath = slashReplaceRegexp.ReplaceAllString(fixImgRelPath, "/")

	return fixImgRelPath, absImgPath, nil
}
