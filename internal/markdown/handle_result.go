package markdown

import (
	"fmt"
	"strings"

	"github.com/bunnier/mdim/internal/base"
)

type MarkdownHandleResult struct {
	DocPath    string
	AllRefImgs base.Set

	RelPathCannotFixedErr   []error
	HasChangeDuringMaintain bool
	SavedMaintainResult     bool

	WebImgDownloadErr []error

	Err error
}

func (handleResult MarkdownHandleResult) ToString() string {
	var resultSb strings.Builder
	switch {
	case handleResult.Err != nil:
		resultSb.WriteString(fmt.Sprintf("[markdown handle]:Occur an error when maintain doc.\n----> %s\n----> %s", handleResult.DocPath, handleResult.Err.Error()))
	case !handleResult.HasChangeDuringMaintain:
		resultSb.WriteString(fmt.Sprintf("[markdown handle]:Correct doc.\n----> %s", handleResult.DocPath))
	case !handleResult.SavedMaintainResult:
		resultSb.WriteString(fmt.Sprintf("[markdown handle]:Find relative path error or web image in a doc, do not change document this time.\n----> %s", handleResult.DocPath))
	case handleResult.SavedMaintainResult:
		resultSb.WriteString(fmt.Sprintf("[markdown handle]:Find relative path error or web image in a doc, change the document successfully.\n----> %s", handleResult.DocPath))
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

	return resultSb.String()
}
