package types

import (
	"fmt"
	"strings"
)

type MarkdownHandleResult struct {
	DocPath string

	RelPathCannotFixedErr []error
	HasErrImgRelPath      bool
	FixedErrImgRelPath    bool

	Err error
}

func (handleResult MarkdownHandleResult) ToString() string {
	var resultSb strings.Builder
	switch {
	case handleResult.Err != nil:
		resultSb.WriteString(fmt.Sprintf("[markdown handle]:Occur an error when maintain doc.\n----> %s\n----> %s", handleResult.DocPath, handleResult.Err.Error()))
	case !handleResult.HasErrImgRelPath:
		resultSb.WriteString(fmt.Sprintf("[markdown handle]:Correct doc.\n----> %s", handleResult.DocPath))
	case !handleResult.FixedErrImgRelPath:
		resultSb.WriteString(fmt.Sprintf("[markdown handle]:Find relative path error in doc, do not fix this time.\n----> %s", handleResult.DocPath))
	case handleResult.FixedErrImgRelPath:
		resultSb.WriteString(fmt.Sprintf("[markdown handle]:Find relative path error in doc, fix successfully.\n----> %s", handleResult.DocPath))
	default:
		resultSb.WriteString("Impossible error.")
	}

	if len(handleResult.RelPathCannotFixedErr) > 0 {
		resultSb.WriteString("\n------> Notice! Have unresolve images.")
		for _, v := range handleResult.RelPathCannotFixedErr {
			resultSb.WriteString(fmt.Sprintf("\n--------> %s", v.Error()))
		}
	}

	return resultSb.String()
}
