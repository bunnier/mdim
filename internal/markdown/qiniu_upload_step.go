package markdown

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bunnier/mdim/internal/qiniu"
)

// NewUploadLocalImgToQiniuStep
func NewUploadLocalImgToQiniuStep(api *qiniu.QiniuUploadApi) ImageMaintainStep {
	return func(imgTag *MarkdownImageTag, handleResult MarkdownHandleResult) error {
		if imgTag.IsWebUrl {
			return nil // No suit.
		}

		if fixedPath, absFixedPath, err := doQiniuUpload(imgTag.DocPath, imgTag.ImgPath, imgTag.AbsImgFolder, api); err != nil {
			handleResult.UploadToQiniuErr = append(handleResult.UploadToQiniuErr, err)
			return err
		} else {
			handleResult.AllRefImgs.Add(absFixedPath)
			imgTag.WholeTag = fmt.Sprintf("![%s](%s)", imgTag.ImgTitle, fixedPath)
			return nil
		}
	}
}

func doQiniuUpload(docPath string, imgPath string, absImgFolder string, api *qiniu.QiniuUploadApi) (string, string, error) {
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

	if url, err := api.Upload("", fixImgAbsPath); err != nil {
		return "", "", err
	} else {
		return url, fixImgAbsPath, nil
	}
}
