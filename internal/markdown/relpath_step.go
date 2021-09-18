package markdown

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// MaintainImageRelateUrl
func MaintainImageRelateUrl(imgTag *MarkdownImageTag, handleResult MarkdownHandleResult) error {
	if fixedPath, absFixedPath, err := getFixImgRelpath(imgTag.DocPath, imgTag.ImgPath, imgTag.AbsImgFolder); err != nil {
		if handleResult.RelPathCannotFixedErr == nil {
			handleResult.RelPathCannotFixedErr = make([]error, 0, 1)
		}
		handleResult.RelPathCannotFixedErr = append(handleResult.RelPathCannotFixedErr, err)
		return err
	} else {
		handleResult.AllRefImgs.Add(absFixedPath)
		imgTag.WholeTag = fmt.Sprintf("![%s](%s)", imgTag.ImgTitle, fixedPath)
		return nil
	}
}

// TryFixLocalImageRelpath
func TryFixLocalImageRelpath(imgTag *MarkdownImageTag, handleResult MarkdownHandleResult) error {
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
		imgTag.WholeTag = fmt.Sprintf("![%s](%s)", imgTag.ImgTitle, fixedPath)
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
