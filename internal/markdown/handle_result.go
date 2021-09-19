package markdown

import (
	"fmt"
	"strings"

	"github.com/bunnier/mdim/internal/base"
)

type MarkdownHandleResult struct {
	DocPath    string
	AllRefImgs base.Set

	HasChangeDuringWorkflow bool
	SavedResult             bool

	RelPathCannotFixedErr []error
	WebImgDownloadErr     []error
	UploadToQiniuErr      []error

	Err error
}

func (handleResult MarkdownHandleResult) String() string {
	var resultSb strings.Builder
	switch {
	case handleResult.Err != nil:
		resultSb.WriteString(fmt.Sprintf("[markdown handle]: Occur an error when maintain doc.\n----> %s\n----> %s", handleResult.DocPath, handleResult.Err.Error()))
	case !handleResult.HasChangeDuringWorkflow:
		resultSb.WriteString(fmt.Sprintf("[markdown handle]: Nothing to do for document.\n----> %s", handleResult.DocPath))
	case !handleResult.SavedResult:
		resultSb.WriteString(fmt.Sprintf("[markdown handle]: Exec workflow successfully, but did not save document this time.\n----> %s", handleResult.DocPath))
	case handleResult.SavedResult:
		resultSb.WriteString(fmt.Sprintf("[markdown handle]: Exec workflow and save the document successfully.\n----> %s", handleResult.DocPath))
	default:
		resultSb.WriteString("Impossible error.")
	}

	if len(handleResult.RelPathCannotFixedErr) > 0 {
		resultSb.WriteString("\n------> Notice! Have unresolve images.")
		for _, v := range handleResult.RelPathCannotFixedErr {
			resultSb.WriteString(fmt.Sprintf("\n--------> %s", v.Error()))
		}
	}

	if len(handleResult.WebImgDownloadErr) > 0 {
		resultSb.WriteString("\n------> Notice! Error occured during downloading images.")
		for _, v := range handleResult.WebImgDownloadErr {
			resultSb.WriteString(fmt.Sprintf("\n--------> %s", v.Error()))
		}
	}

	if len(handleResult.UploadToQiniuErr) > 0 {
		resultSb.WriteString("\n------> Notice! Error occured during uploading images.")
		for _, v := range handleResult.UploadToQiniuErr {
			resultSb.WriteString(fmt.Sprintf("\n--------> %s", v.Error()))
		}
	}

	return resultSb.String()
}
