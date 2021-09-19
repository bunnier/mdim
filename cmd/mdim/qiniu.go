package main

import (
	"fmt"
	"os"

	"github.com/bunnier/mdim/internal/base"
	"github.com/bunnier/mdim/internal/markdown"
	"github.com/bunnier/mdim/internal/qiniu"
	"github.com/spf13/cobra"
)

// QiniuOptions are command-line options.
type QiniuOptions struct {
	Sk       string
	Ak       string
	Bucket   string
	Folder   string
	Domain   string
	UseHTTPS bool
}

var qiniuOptions = &QiniuOptions{}

func init() {
	flags := qiniuCmd.Flags()
	initBaseOptions(flags)
	flags.StringVarP(&qiniuOptions.Sk, "sk", "", "", "Must not be empty. Assign the SK(Secret Key) of Qiniu SDK, also can be provided by setting env variable named 'mdim_qiniu_sk'.")
	flags.StringVarP(&qiniuOptions.Ak, "ak", "", "", "Must not be empty. Assign the AK(Access Key) of Qiniu SDK, also can be provided by setting env variable named 'mdim_qiniu_ak'.")
	flags.StringVarP(&qiniuOptions.Bucket, "bucket", "b", "", "Must not be empty. Assign the Bucket of Qiniu SDK, also can be provided by setting env variable named 'mdim_qiniu_bucket'.")
	flags.StringVarP(&qiniuOptions.Domain, "domain", "d", "", "The domain of uploaded image url, also can be provided by setting env variable named 'mdim_qiniu_domain'. If do not assign the option, will use first domain in specific bucket")
	flags.StringVarP(&qiniuOptions.Folder, "folder", "u", "", "After uploaded, you images can access in this url {protocal}://{domain}/{path}/{img_name}")
	flags.BoolVarP(&qiniuOptions.UseHTTPS, "https", "t", false, "If assign this option, will use https instead of http.")
}

var qiniuCmd = &cobra.Command{
	Use:   "qiniu",
	Short: "Uploading the local image files in specific markdown files to Qiniu cloud space.",
	Long:  "The tool helps to upload the local image files in specific markdown files to Qiniu cloud space.",
	Run: func(cmd *cobra.Command, args []string) {
		validateBaseOptions(cmd)

		if qiniuOptions.Ak == "" {
			qiniuOptions.Ak = os.Getenv("mdim_qiniu_ak")
		}

		if qiniuOptions.Sk == "" {
			qiniuOptions.Sk = os.Getenv("mdim_qiniu_sk")
		}

		if qiniuOptions.Bucket == "" {
			qiniuOptions.Bucket = os.Getenv("mdim_qiniu_bucket")
		}

		if qiniuOptions.Domain == "" {
			qiniuOptions.Domain = os.Getenv("mdim_qiniu_domain")
		}

		if qiniuOptions.Ak == "" || qiniuOptions.Sk == "" || qiniuOptions.Bucket == "" {
			cmd.Usage()
			os.Exit(1)
		}

		doQiniuCmd(qiniuOptions)
	},
}

func doQiniuCmd(param *QiniuOptions) {
	fmt.Println("==========================================")
	fmt.Println("  Starting to scan markdown document(s)..")
	fmt.Println("==========================================")

	qiniuApi := qiniu.NewQuniuUploadApi(
		param.Ak, param.Sk, param.Bucket, param.UseHTTPS,
		qiniu.QiniuUploadApiDomainOption(param.Domain),
		qiniu.QiniuUploadApiDefaultFolderOption(param.Folder),
	)

	// Scan docs in docFolder to maintain image tags.
	markdownHandleResults := markdown.WalkDirToHandleDocs(
		baseOptions.SingleDocument, baseOptions.AbsDocFolder, baseOptions.AbsImgFolder, baseOptions.DoSave,
		[]markdown.ImageMaintainStep{markdown.NewUploadLocalImgToQiniuStep(qiniuApi)})

	hasInterruptErr := false
	allRefImgsAbsPathSet := base.NewSet(100)
	for _, handleResult := range markdownHandleResults {
		if handleResult.HasChangeDuringWorkflow || handleResult.UploadToQiniuErr != nil {
			fmt.Println(handleResult.String())
			fmt.Println()
		}

		if handleResult.Err != nil {
			hasInterruptErr = true
		}

		allRefImgsAbsPathSet.Merge(handleResult.AllRefImgs)
	}

	if hasInterruptErr {
		os.Exit(10)
	}

	fmt.Println("==========================================")
	fmt.Println("  All done.")
	fmt.Println("==========================================")
}
