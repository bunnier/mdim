package types

type MarkdownDocHandleResult struct {
	DocPath            string
	HasErrImgRelPath   bool
	FixedErrImgRelPath bool
	Err                error
}

type ImageHandleResult struct {
	ImagePath string
	Deleted   bool
	Err       error
}
