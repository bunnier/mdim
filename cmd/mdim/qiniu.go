package main

import (
	"os"

	"github.com/spf13/cobra"
)

// QiniuParam are command-line options.
type QiniuParam struct {
	Sk     string
	Ak     string
	Bucket string
}

var qiniuParam = &QiniuParam{}

func init() {
	flags := qiniuCmd.Flags()
	flags.StringVarP(&qiniuParam.Sk, "sk", "s", "", "Must not be empty. Assign the SK(Secret Key) of Qiniu SDK.")
	flags.StringVarP(&qiniuParam.Ak, "ak", "a", "", "Must not be empty. Assign the AK(Access Key) of Qiniu SDK.")
	flags.StringVarP(&qiniuParam.Bucket, "bucket", "b", "", "Must not be empty. Assign the Bucket of Qiniu SDK.")
}

var qiniuCmd = &cobra.Command{
	Use:   "qiniu",
	Short: "Uploading the local image files in specific markdown files to Qiniu cloud space.",
	Long:  "The tool helps to upload the local image files in specific markdown files to Qiniu cloud space.",
	Run: func(cmd *cobra.Command, args []string) {
		if qiniuParam.Ak == "" {
			qiniuParam.Ak = os.Getenv("mdim_qiniu_ak")
		}

		if qiniuParam.Sk == "" {
			qiniuParam.Ak = os.Getenv("mdim_qiniu_sk")
		}

		if qiniuParam.Bucket == "" {
			qiniuParam.Ak = os.Getenv("mdim_qiniu_bucket")
		}

		if qiniuParam.Ak == "" || qiniuParam.Sk == "" || qiniuParam.Bucket == "" {
			cmd.Usage()
			os.Exit(1)
		}

		doQiniuCmd(qiniuParam)
	},
}

func doQiniuCmd(param *QiniuParam) {
}
