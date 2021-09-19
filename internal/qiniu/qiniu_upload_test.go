package qiniu

import (
	"net/url"
	"path/filepath"
	"strings"
	"testing"
)

const (
	ak          string = ""
	sk          string = ""
	bucket      string = ""
	testImgPath string = ""
	testDomain  string = ""
	useHTTPS    bool   = true
)

func TestQiniuUploadApi_Upload(t *testing.T) {
	type args struct {
		remoteFilepath string
		localFilepath  string
	}
	tests := []struct {
		name    string
		api     *QiniuUploadApi
		args    args
		wantErr bool
	}{
		{
			"test1",
			NewQuniuUploadApi(ak, sk, bucket, useHTTPS, QiniuUploadApiDomainOption(testDomain)),
			args{remoteFilepath: filepath.Base(testImgPath), localFilepath: testImgPath},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantUrl := strings.Join([]string{tt.api.domain, url.PathEscape(tt.args.remoteFilepath)}, "/")
			var url string
			var err error
			if url, err = tt.api.Upload(tt.args.remoteFilepath, tt.args.localFilepath); (err != nil) != tt.wantErr {
				t.Errorf("QiniuUploadApi.Upload() error = %v, wantErr %v", err, tt.wantErr)
			}

			if wantUrl != url {
				t.Errorf("QiniuUploadApi.Upload() url = %v, wantUrl %v", url, wantUrl)
			}
		})
	}
}
