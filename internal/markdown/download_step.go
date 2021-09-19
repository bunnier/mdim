package markdown

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var (
	// Input=Relative path, Group1=First named folder name, Group2=relative path in imgFolder
	imgPathRegexp *regexp.Regexp = regexp.MustCompile(`^(?:(?:\.{1,2}[/\\])+)([^/\\\n]+)?[/\\](.+)$`)

	// For get web image suffix.
	httpImgRegex = regexp.MustCompile(`^(?:http[s]?://)(?:[^/]+/)+.+(\.[a-zA-Z]{1,5})$`)

	// For slash replace.
	slashReplaceRegexp = regexp.MustCompile(`\\`)
)

// DownloadImageStep
func DownloadImageStep(imgTag *MarkdownImageTag, handleResult MarkdownHandleResult) error {
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
		imgTag.WholeTag = fmt.Sprintf("![%s](%s)", imgTag.ImgTitle, fixedPath)
		return nil
	}
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
