package internal

import (
	"path/filepath"
	"testing"
)

const (
	ak          string = ""
	sk          string = ""
	bucket      string = ""
	testImgPath string = ""
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
			NewQuniuUploadApi(ak, sk, bucket),
			args{remoteFilepath: filepath.Base(testImgPath), localFilepath: testImgPath},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.api.Upload(tt.args.remoteFilepath, tt.args.localFilepath); (err != nil) != tt.wantErr {
				t.Errorf("QiniuUploadApi.Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQiniuUploadApi_GetDomains(t *testing.T) {
	tests := []struct {
		name    string
		api     *QiniuUploadApi
		wantErr bool
	}{
		{
			"test1",
			NewQuniuUploadApi(ak, sk, bucket),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.api.GetDomains(); (err != nil) != tt.wantErr {
				t.Errorf("QiniuUploadApi.GetDomains() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
